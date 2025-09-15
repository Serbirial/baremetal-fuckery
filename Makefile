
# TinyGo compiler command
TINYGO = tinygo

# Target JSON for Pi Zero 2W
TARGET = ./rpi_zero2w.json

# ARM cross-linker (for objcopy to raw binary)
OBJCOPY = aarch64-linux-gnu-objcopy

# Output filenames
ELF = build.elf
IMG = boot/kernel8.img

# All steps: compile Go code, link with boot.S, produce binary image
all: $(IMG)

# -------------------------------
# Step 1: Build ELF with TinyGo
# -------------------------------
$(ELF): sys/init.go sys/uart.go boot/boot.S boot/linker_rpi.ld $(TARGET)
	@echo "=== Building ELF with TinyGo ==="
	$(TINYGO) build -o $(ELF) -target=$(TARGET) ./sys

# -------------------------------
# Step 2: Convert ELF to raw binary
# -------------------------------
$(IMG): $(ELF)
	@echo "=== Converting ELF to kernel8.img ==="
	$(OBJCOPY) -O binary $(ELF) $(IMG)

# -------------------------------
# Clean build artifacts
# -------------------------------
clean:
	@echo "=== Cleaning build artifacts ==="
	rm -f $(ELF) $(IMG)

# -------------------------------
# Flash instructions (manual)
# -------------------------------
flash:
	@echo "=== Copy $(IMG) to SD card boot partition ==="
	@echo "Use your OS file manager or 'sudo cp $(IMG) /media/pi/boot/'"