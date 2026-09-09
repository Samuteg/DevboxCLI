package scaffold

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// maxTemplateBytes limita o zip lido (anti zip-bomb). Templates atuais têm
// ~81M (node); 128M dá margem sem aceitar arquivos absurdos.
const maxTemplateBytes = 128 << 20

func MaterializeTemplates(fsys fs.FS, sourceZip, targetRoot string) error {
	f, err := fsys.Open(sourceZip)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := io.ReadAll(io.LimitReader(f, maxTemplateBytes+1))
	if err != nil {
		return err
	}
	if int64(len(data)) > maxTemplateBytes {
		return fmt.Errorf("template %s excede o limite de %d bytes", sourceZip, maxTemplateBytes)
	}

	zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return err
	}

	absRoot, err := filepath.Abs(targetRoot)
	if err != nil {
		return err
	}

	for _, file := range zipReader.File {
		cleanName := filepath.Clean(file.Name)
		if cleanName == "." || strings.HasPrefix(cleanName, "..") || filepath.IsAbs(cleanName) {
			return fmt.Errorf("caminho inseguro no zip: %s", file.Name)
		}
		// Rejeita symlinks (podem apontar para fora da árvore).
		if file.Mode()&fs.ModeSymlink != 0 {
			return fmt.Errorf("symlink não permitido no zip: %s", file.Name)
		}
		targetPath := filepath.Join(targetRoot, cleanName)
		absTarget, err := filepath.Abs(targetPath)
		if err != nil {
			return err
		}
		if absTarget != absRoot && !strings.HasPrefix(absTarget, absRoot+string(filepath.Separator)) {
			return fmt.Errorf("caminho escapa do destino: %s", file.Name)
		}

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return err
		}

		rc, err := file.Open()
		if err != nil {
			return err
		}

		// Máscara conservadora: nunca suid/sgid/sticky, sem exec salvo .sh.
		mode := file.Mode().Perm()
		if mode == 0 {
			mode = 0644
		}
		mode &= 0777
		mode &^= 07000
		if filepath.Ext(cleanName) != ".sh" {
			mode &^= 0111
		}
		outFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
		if err != nil {
			rc.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		rc.Close()
		outFile.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
