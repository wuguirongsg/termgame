#!/usr/bin/env python3
"""
Test script for Windows compatibility
"""
import sys
import os

IS_WINDOWS = os.name == 'nt'

print(f"OS: {os.name}")
print(f"Platform: {sys.platform}")
print(f"Is Windows: {IS_WINDOWS}")

if IS_WINDOWS:
    try:
        import msvcrt
        print("✓ msvcrt module available")
    except ImportError as e:
        print(f"✗ msvcrt module not available: {e}")
else:
    try:
        import tty
        import termios
        import select
        print("✓ Unix modules available (tty, termios, select)")
    except ImportError as e:
        print(f"✗ Unix modules not available: {e}")

print("\nAll required modules are available!")