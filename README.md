<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Devbox CLI - Sua caixa de ferramentas no terminal</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;600;700&family=JetBrains+Mono:wght@400;700&display=swap" rel="stylesheet">
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css">

    <script>
        tailwind.config = {
            theme: {
                extend: {
                    colors: {
                        go: '#00ADD8',
                        dark: '#1a1b26',
                        darker: '#16161e',
                        card: '#24283b',
                        accent: '#bb9af7'
                    },
                    fontFamily: {
                        sans: ['Inter', 'sans-serif'],
                        mono: ['JetBrains Mono', 'monospace'],
                    }
                }
            }
        }
    </script>
    <style>
        body { background-color: #1a1b26; color: #a9b1d6; }
        .gradient-text {
            background: linear-gradient(to right, #00ADD8, #bb9af7);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
        }
        .code-block {
            background-color: #16161e;
            border: 1px solid #414868;
        }
    </style>
</head>
<body class="antialiased selection:bg-go selection:text-white">

    <header class="max-w-5xl mx-auto px-6 py-20 text-center">
        <div class="mb-6 flex justify-center gap-2">
            <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go" alt="Go Version">
            <img src="https://img.shields.io/badge/Platform-Win%20%7C%20Linux-brightgreen" alt="Platform">
            <img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="License">
        </div>
        
        <h1 class="text-6xl font-bold mb-4 text-white tracking-tight">
            📦 <span class="gradient-text">Devbox CLI</span>
        </h1>
        <p class="text-xl mb-10 max-w-2xl mx-auto text-gray-400">
            Acelere seu fluxo de trabalho combinando a robustez do <span class="text-go font-bold">Go</span> com uma interface interativa moderna.
        </p>

        <div class="flex justify-center gap-4">
            <a href="#install" class="bg-go hover:bg-cyan-600 text-white font-bold py-3 px-8 rounded-full transition transform hover:scale-105 shadow-lg shadow-cyan-500/20">
                <i class="fas fa-download mr-2"></i> Instalar Agora
            </a>
            <a href="https://github.com/seu-usuario/devbox" target="_blank" class="bg-card hover:bg-gray-700 text-white font-bold py-3 px-8 rounded-full transition border border-gray-600">
                <i class="fab fa
