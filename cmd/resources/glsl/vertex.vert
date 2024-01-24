#version 440

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

in vec3 vert;

uniform vec4 diffuse;
out vec4 vertexColor;

void main() {
    vertexColor = diffuse;
    gl_Position = projection * camera * model * vec4(vert, 1);
}
