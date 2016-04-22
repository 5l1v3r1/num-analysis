rm -rf build/
mkdir build
cd src
gopherjs build -m
mv src.js ../build/src.js
mv src.js.map ../build/src.js.map
