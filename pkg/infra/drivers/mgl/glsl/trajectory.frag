#version 430 core

in vec4 color4;
in vec2 markerCoordinate;
flat in float isMarker;
out vec4 outColor;

void main() {
    if (isMarker > 0.5 && dot(markerCoordinate, markerCoordinate) > 1.0) {
        discard;
    }
    outColor = color4;
}
