"use strict";

const pixels = new Uint8ClampedArray(8192);

let gameInterval;
let wasm;

const canvas = document.getElementById("screen");
const ctx = canvas.getContext("2d");

document.onkeydown = obj => keyEvent(obj, true);
document.onkeyup = obj => keyEvent(obj, false);
const keys = [88, 49, 50, 51, 81, 87, 69, 65, 83, 68, 90, 67, 52, 82, 70, 86];

const roms = ["15PUZZLE", "BLINKY", "BLITZ", "BRIX", "CONNECT4", "GUESS", "HIDDEN", "INVADERS", "KALEID", "MAZE", "MERLIN", "MISSLE", "PONG", "PONG2", "PUZZLE", "SYZYGY", "TANK", "TETRIS", "TICTAC", "UFO", "VBRIX", "VERS", "WIPEOFF"];
const fileElement = document.createElement("input");
fileElement.setAttribute("type", "file");
fileElement.onchange = fileRom;

const AudioContext = window.AudioContext || window.webkitAudioContext;
const audio = new AudioContext();

let oscillator;
const gainNode = audio.createGain();
gainNode.gain.setValueAtTime(0, audio.currentTime);

// var is needed in order for WebAssembly to be able to find the value in global
var rom;

function draw(x, y, data) {
    if (y >= 32) return false;
    const index = 256*y + 4*x;
    let result = false;
    for (let i = 0; i < 8; i++) {
        const bit = (data >> (7-i)) & 1;
        const value = bit == 1 ? 255 : 0;
        if (bit == 1 && pixels[index+4*i] != 0)
            result = true;
        for (let j = 0; j < 3; j++) {
            pixels[index+4*i+j] ^= value;
        }
    }
    return result;
}

function clear() {
    for (let i = 0; i < pixels.length; i++) {
        if (i%4 == 3) {
            pixels[i] = 255;
        } else {
            pixels[i] = 0;
        }
    }
}

function pause() {
    if (oscillator) oscillator.stop();
    window.clearInterval(gameInterval);
    gainNode.gain.setValueAtTime(0, audio.currentTime);
}

function runFrame() {
    const st = wasm.exports.runFrame();
    ctx.putImageData(new ImageData(pixels, 64, 32), 0, 0);

    const gain = st ? 0.1 : 0;
    gainNode.gain.setValueAtTime(gain, audio.currentTime);
}

function play() {
    oscillator = audio.createOscillator();
    oscillator.type = "square";
    oscillator.frequency.setValueAtTime(440, audio.currentTime);
    oscillator.connect(gainNode).connect(audio.destination);
    oscillator.start();

    window.clearInterval(gameInterval);
    gameInterval = window.setInterval(runFrame, 1000/60);
}

function reset() {
    pause();
    wasm.exports.loadRom();
    document.getElementById("gameList").className = "hidden";
    document.getElementById("game").className = "";
    play();
}

function loadRom(romPath) {
    window.clearInterval(gameInterval);

    fetch("../chip8-games/" + romPath).then(
        resp => resp.status == 200 ? resp.body.getReader() : undefined
    ).then(
        reader => reader.read()
    ).then(obj => {
        rom = obj.value;
        reset();
    });
}

function fileRom() {
    window.clearInterval(gameInterval);

    fileElement.files[0].stream().getReader().read().then(obj => {
        rom = obj.value;
        reset();
    });
}

function chooseRom() {
    pause();
    document.getElementById("game").className = "hidden";
    document.getElementById("gameList").className = "";
}

function keyEvent(obj, down) {
    if (!wasm) return;
    wasm.exports.setKey(keys.indexOf(obj.keyCode), down);
}

function init() {
    for (let i = 3; i < pixels.length; i += 4) {
        pixels[i] = 255;
    }

    const go = new Go();
    WebAssembly.instantiateStreaming(fetch("chip8-tiny.wasm"), go.importObject).then(obj => {
        wasm = obj.instance;
        go.run(wasm);
        loadRom("INVADERS");
    });

    clear();
    ctx.putImageData(new ImageData(pixels, 64, 32), 0, 0);

    const games = document.getElementById("gameList");
    const load = document.createElement("input");
    load.type = "button";
    load.onclick = () => fileElement.click();
    load.value = "Load ROM from your computer";
    games.appendChild(load);
    for (const rom of roms) {
        const element = document.createElement("input");
        element.type = "button";
        element.onclick = () => loadRom(rom);
        element.value = rom;
        games.appendChild(document.createElement("br"));
        games.appendChild(element);
    }
}

init();
