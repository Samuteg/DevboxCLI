package scaffold

import "path/filepath"

type Variant struct {
	Name      string
	Source    string
	ExtraDirs []string
}

type Stack struct {
	Name       string
	IsBackend  bool
	Source     string
	ExtraDirs  []string
	RunInstall bool
	Variants   []Variant
}

func DefaultStacks() map[string]Stack {
	pythonSubDirs := []string{
		"src/api/routes",
		"src/api/dependencies",
		"src/config/settings",
		"src/models/schemas",
		"src/models/db",
		"src/repository/crud",
		"src/repository/migrations/versions",
		"src/securities/authorizations",
		"src/securities/hashing",
		"src/securities/verifications",
		"src/utilities/exceptions/http",
		"src/utilities/formatters",
		"src/utilities/messages/exceptions/http",
		"tests/end_to_end_tests",
		"tests/unit_tests",
		"tests/integration_tests",
		"tests/security_tests",
	}

	golangSubDirs := []string{
		"cmd/apps",
		"docs/images",
		"pkg/conf",
		"pkg/logger",
		"pkg/server/auth/acl",
		"pkg/server/controller",
		"pkg/server/model",
		"pkg/server/router/middleware",
		"pkg/server/storage/cache/local",
		"pkg/server/storage/cache/redis",
		"pkg/server/storage/db/mysql",
	}

	rubySubDirs := []string{
		".github/workflows",
		"lib/generators/rails/templates",
		"lib/jb",
		"test/dummy_app/app/assets/config",
		"test/dummy_app/app/assets/javascripts",
		"test/dummy_app/app/assets/stylesheets",
		"test/dummy_app/app/controllers",
		"test/dummy_app/app/helpers",
		"test/dummy_app/app/mailers",
		"test/dummy_app/app/models",
		"test/dummy_app/app/views",
		"test/dummy_app/bin",
		"test/dummy_app/config/environments",
		"test/dummy_app/config/initializers",
		"test/dummy_app/config/features",
	}

	return map[string]Stack{
		"Node.js": {
			Name:      "Node",
			IsBackend: true,
			Source:    "templates/node/base",
			Variants: []Variant{
				{
					Name:      "JavaScript",
					Source:    "templates/node/js",
					ExtraDirs: []string{"src/controllers", "src/libs", "src/middleware", "src/routes", "src/models"},
				},
				{
					Name:      "TypeScript",
					Source:    "templates/node/ts",
					ExtraDirs: []string{"src/controllers", "src/libs", "src/middleware", "src/routes", "src/models"},
				},
			},
			RunInstall: true,
		},
		"Go": {
			Name:      "Go",
			IsBackend: true,
			Source:    "templates/go",
			Variants: []Variant{
				{Name: "Simples (padrao)", Source: "templates/golang/simple", ExtraDirs: []string{"cmd/api", "internal/entity", "internal/infra/repository", "internal/infra/web", "internal/usecase"}},
				{Name: "Gin", Source: "templates/golang/Gin", ExtraDirs: golangSubDirs},
			},
		},
		"Python": {
			Name:      "Python",
			IsBackend: true,
			Source:    "templates/python",
			Variants: []Variant{
				{Name: "Simples (Recomendado)", Source: "templates/python/simple"},
				{Name: "FastAPI", Source: "templates/python/FastAPI", ExtraDirs: prefixPaths("backend", pythonSubDirs)},
			},
		},
		"Ruby": {
			Name:      "Ruby",
			IsBackend: true,
			Source:    "templates/ruby",
			Variants: []Variant{
				{Name: "simple", Source: "templates/ruby/simple"},
				{Name: "Rails", Source: "templates/ruby/On_rails", ExtraDirs: rubySubDirs},
			},
		},
		"Next.js": {Name: "Next", IsBackend: false, Source: "pnpm create next-app@latest %s"},
		"Vite":    {Name: "Vite", IsBackend: false, Source: "pnpm create vite@latest %s"},
	}
}

func prefixPaths(prefix string, subs []string) []string {
	fullPaths := make([]string, len(subs))
	for i, sub := range subs {
		fullPaths[i] = filepath.Join(prefix, sub)
	}
	return fullPaths
}
