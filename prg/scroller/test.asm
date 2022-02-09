
.const debug = false
:BasicUpstart2(start)

start: {
    lda $DD00 // 081C
    and #%11111100
    ora #%00000010 // VIC-II bank 2
    sta $dd00

    lda #0
    sta $d020
    sta $d021
    lda #$18
    sta $d018

    // clear out two last rows of the bitmap
    lda #0 //0833
    .for (var i = 0; i < 40; i++) {
        sta colora + 40*23 + i
        sta colorb + 40*23 + i
    }

infloop:
    jmp infloop
}

.align 64
* = $6000 "image"
gfx:
.import source "titleimage.txt"