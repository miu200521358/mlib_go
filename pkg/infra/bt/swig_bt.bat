swig -c++ -go -cgo -I"C:\MMD\mlib_go\pkg\infra\bt" ^
    -I"C:\development\TDM-GCC-64\lib\gcc\x86_64-w64-mingw32\10.3.0\include\c++\x86_64-w64-mingw32" ^
    -I"C:\development\TDM-GCC-64\x86_64-w64-mingw32\include" ^
    -I"C:\Program Files\Microsoft Visual Studio\2022\Community\VC\Tools\MSVC\14.38.33130\include" ^
    -cpperraswarn -o "C:\MMD\mlib_go\pkg\infra\bt\bt.cxx" "C:\MMD\mlib_go\pkg\infra\bt\bullet.i"

go clean --cache
