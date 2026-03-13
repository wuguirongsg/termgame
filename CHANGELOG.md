# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial release of FC Terminal Mario
- Java implementation with Maven build system
- Python implementation
- Complete level with platforms, pipes, enemies, and coins
- Physics engine with inertia-based movement
- Collision detection system
- Animation system for Mario and Goombas
- Score and coin collection
- Lives system
- Game over and victory states
- HUD display (score, coins, lives, time)

### Features
- Smooth acceleration and deceleration for character movement
- Air control (reduced control while jumping)
- Multi-key support (can press left and right simultaneously)
- Question blocks that spawn coins
- Goomba enemies with AI
- Staircase approach to flagpole
- Cloud and hill decorations

### Controls
- A / Left Arrow - Move left
- D / Right Arrow - Move right
- Space / W / Up Arrow - Jump
- Q - Quit game

### Technical
- ANSI escape sequence rendering
- Half-block character (▀) for double vertical resolution
- Raw terminal mode for non-blocking input
- 60 FPS game loop
- Cross-platform support (macOS, Linux, Windows)

### Fixed
- Windows console input handling (use JLine library)
- Windows console output encoding (UTF-8 with chcp 65001)
- Removed emoji characters that caused encoding issues on Windows
- Python version Windows compatibility (msvcrt for input)
- Arrow key detection on Python version

## [1.0.0] - 2026-03-13

### Added
- Initial public release
- Java and Python implementations
- Complete gameplay mechanics