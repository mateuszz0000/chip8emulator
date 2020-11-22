package chip8

import "github.com/hajimehoshi/ebiten"
import "image/color"
import "fmt"


type Chip8 struct {
	Memory [4096]byte
	ProgramCounter uint16
	Registers [16]byte
	IndexRegister uint16
	Stack [16]uint16
	StackPointer byte 
	Screen *ebiten.Image
	Timer byte
	SoundTimer byte
	Pixels [64][32]byte
	Keyboard [16]byte
	ClearScreenFlag bool
}

var Emulator (*Chip8) // global pointer to struct

var square, _ = ebiten.NewImage(10, 10, ebiten.FilterNearest)

func (chip8 *Chip8) Run() {
	square.Fill(color.White)
	// ebiten.SetMaxTPS(120)
    if err := ebiten.Run(update, 640, 320, 1, "Hello world!"); err != nil {
        panic(err)
    }
}

func update(screen *ebiten.Image) error {
	screen.Fill(color.NRGBA{0x00, 0x00, 0x00, 0xff})
	for elo :=0; elo <15; elo++ {
		Emulator.runCycle()
		Emulator.checkClearScreen(screen)
		Emulator.render(screen)
		CheckKeyboard(Emulator)
}
	return nil
}

func (chip8 *Chip8) runCycle(){
	opcode := chip8.getOpcode()

	fmt.Println("fetched op", fmt.Sprintf("%X", opcode))
	chip8.ProgramCounter += 2
	chip8.decodeAndRunInstruction(opcode)

	if chip8.Timer > 0 {
		chip8.Timer--
	}

	if chip8.SoundTimer > 0 {
		chip8.SoundTimer--
		// TODO: beep buzzer
	}
}

func (chip8 *Chip8) render(screen *ebiten.Image) {
	for x:=0 ; x < 64; x++ {
		for y:=0 ; y< 32; y++{
			if chip8.Pixels[x][y] == 0xFF{ 
				
				opts := &ebiten.DrawImageOptions{}
				opts.GeoM.Translate(float64(x * 10), float64(y*10))
				screen.DrawImage(square, opts)
			}
			
		}
	}
}

func (chip8 *Chip8) checkClearScreen(screen *ebiten.Image) {
	if chip8.ClearScreenFlag{
		for x:=0 ; x < 64; x++ {
			for y:=0 ; y< 32; y++{
				chip8.Pixels[x][y] = 0x00
			}
		}

		screen.Fill(color.NRGBA{0x00, 0x00, 0x00, 0xff})
		chip8.ClearScreenFlag = false
	}
	
}

func (chip8 *Chip8) getOpcode() uint16{
	addr := chip8.ProgramCounter
	return uint16(chip8.Memory[addr]) << 8 | uint16(chip8.Memory[addr + 1])
}


func (chip8 *Chip8) decodeAndRunInstruction(opcode uint16){
	q := opcode >> 12
	opcodes[q](chip8, opcode)

}

