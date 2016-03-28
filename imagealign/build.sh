rm -f assets/script.js
rm -f assets/script.js.map
cd src && gopherjs build -m -o ../assets/script.js
