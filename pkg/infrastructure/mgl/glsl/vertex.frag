#version 440 core

uniform vec4 color;

in float vertexUvX;

out vec4  outColor;

void main() {
    if (vertexUvX < 0) {
        // UVのXが明示的にマイナスの場合、描画しない
        discard;
    }

    outColor = color;
}