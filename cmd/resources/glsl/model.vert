#version 440

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

in layout(location = 0) vec3 position;
in layout(location = 1) vec3 normal;

uniform vec4 diffuse;
uniform vec3 ambient;
uniform vec4 specular;

out vec4 vertexColor;

void main() {
    vertexColor = diffuse;
    gl_Position = projection * camera * model * vec4(position, 1);
}
