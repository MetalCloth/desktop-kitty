# 🐈 Desktop Kitty (GoCat)

![Go Version](https://img.shields.io/github/go-mod/go-version/MetalCloth/desktop-kitty?style=flat&logo=go)
![Platform](https://img.shields.io/badge/Platform-Windows-blue)
![License](https://img.shields.io/badge/License-MIT-green)

> **A chaotic, physics-enabled cat that lives on top of your Windows taskbar.**
> Built from scratch using **Go** and **Ebitengine**.

![Demo Gif](https://media.giphy.com/media/placeholder-url-replace-this-with-your-gif.gif)
*(Replace this line with a GIF of you yeeting the cat!)*

---

## 🚀 Features

* **👻 Transparent Overlay:** Uses system calls to render a borderless, click-through window that sits on top of all other applications.
* **🧠 State Machine AI:** The cat makes decisions based on player interaction—it chases the cursor, idles, sleeps, or demands pets.
* **🚀 "Yeet" Physics:** Implements vector-based momentum transfer. If you drag and release the cat, it calculates your mouse velocity and applies gravity, friction, and wall-bouncing physics.
* **💖 Petting Engine:** Detects "rubbing" motion (mouse velocity + tight hitbox dwell time) to trigger happy states.
* **📦 Single Binary:** Uses Go's `embed` system to bake all assets into a single, portable `.exe` file. No installer or asset folders required.

## 🎮 Controls

| Action | Result |
| :--- | :--- |
| **Move Mouse** | The cat will chase you (unless resting). |
| **Hover** | The cat stops chasing and waits. |
| **Swipe/Rub** | Move mouse back & forth over the cat to **Pet** it (Happy Mode). |
| **Click & Drag** | Pick up the cat (Sticky Hold). |
| **Release** | **YEET** the cat. It will bounce off your screen edges. |

---

## 📥 Installation

**No coding required.**

1.  Go to the **[Releases](../../releases)** page.
2.  Download the latest `SuperCat.exe` (or `GoCat.exe`).
3.  Double-click to run.
4.  *To quit: Go to Task Manager or press Alt+F4 while the cat is active.*

---

## 🛠️ Building from Source

If you are a developer and want to modify the physics or sprites:

### Prerequisites
* [Go 1.22+](https://go.dev/dl/)

### Steps
1.  **Clone the repository:**
    ```bash
    git clone [https://github.com/MetalCloth/desktop-kitty.git](https://github.com/MetalCloth/desktop-kitty.git)
    cd desktop-kitty
    ```

2.  **Install Ebitengine dependencies:**
    ```bash
    go mod tidy
    ```

3.  **Run in Debug Mode:**
    ```bash
    go run .
    ```

4.  **Compile for Windows (Hidden Console):**
    This command builds a standalone `.exe` and hides the terminal window.
    ```bash
    go build -ldflags "-H=windowsgui" -o SuperCat.exe
    ```

---

## 📚 Technical Details

This project explores **2D Vector Math** and **Game Loop Architecture** in Go.

### The Physics Engine
The "Yeet" mechanic relies on calculating velocity (`vx`, `vy`) at the moment of mouse release:
```go
// Calculating momentum on release
g.vx = float64(mx - g.lastMx) * 1.5
g.vy = float64(my - g.lastMy) * 1.5

// Applying gravity and friction every frame
if airborne {
    g.vy += 0.5  // Gravity
    g.vx *= 0.95 // Air Resistance
}

```
### The "Wall Grind" Fix
Standard collision logic often causes sprites to "vibrate" against walls. This project uses **Hysteresis** (different start/stop thresholds) to ensure the cat stops smoothly at the screen border without jittering.

---

## 📜 Credits

* **Code:** [MetalCloth](https://github.com/MetalCloth)
* **Engine:** [Ebitengine](https://ebitengine.org/)
* **Art:** RANDOM PINTEREST SPRITE

---

*Built with ❤️ (and physics) in Golang.*

