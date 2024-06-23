#version 440 core

uniform mat4 modelViewProjectionMatrix;
uniform mat4 modelViewMatrix;

in layout(location = 0) vec3 position;

void main() {
    gl_Position = modelViewProjectionMatrix * modelViewMatrix * vec4(position, 1.0);
}
