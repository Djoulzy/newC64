.macpack cbm
;===============================================================================
.segment "ZEROPAGE"
; $92-$96 (only if no datasette is used)
; $A3-$B1 (only if no RS-232 and datasette is used)
; $F7-$FA (only if no RS-232 is used)
; $FB-$FE (always)

;===============================================================================
; Music include

;===============================================================================
; Sprite include

;===============================================================================
.segment "CODE"
        JMP start               ; run the init code then flow into the update code

;===============================================================================

; Code Includes

;=============================================================================== 

;============================================================
; MAIN
;============================================================

start:
        SEI                  ; set interrupt bit, make the CPU ignore interrupt requests
        LDA #%01111111       ; switch off interrupt signals from CIA-1
        STA $DC0D

        AND $D011            ; clear most significant bit of VIC's raster register
        STA $D011

        LDA $DC0D            ; acknowledge pending interrupts from CIA-1
        LDA $DD0D            ; acknowledge pending interrupts from CIA-2

        LDA #210             ; set rasterline where interrupt shall occur
        STA $D012

        LDA #<Irq            ; set interrupt vectors, pointing to interrupt service routine below
        STA $0314
        LDA #>Irq
        STA $0315

        LDA #%00000001       ; enable raster interrupt signals from VIC
        STA $D01A

        CLI                  ; clear interrupt flag, allowing the CPU to respond to interrupt requests
        RTS

Irq:        
        LDA #$7
        STA $D020            ; change border colour to yellow

        LDX #$90             ; empty loop to do nothing for just under half a millisecond

Pause:     
        DEX
        BNE Pause

        LDA #$0
        STA $D020            ; change border colour to black

        ASL $D019            ; acknowledge the interrupt by clearing the VIC's interrupt flag

        JMP $EA31            ; jump into KERNAL's standard interrupt service routine to handle keyboard scan, cursor display etc.