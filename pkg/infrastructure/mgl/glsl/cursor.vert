#version 440 core

uniform mat4 projectionMatrix;
uniform mat4 viewMatrix;

in layout(location = 0) vec3 position;

void main() {
    gl_Position = projectionMatrix * viewMatrix * vec4(position, 1.0);
}
