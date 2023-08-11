Remove-Item *.db
Remove-Item *.exe
Remove-Item *.dat
Remove-Item *.log

go build -o bc.exe .

./bc.exe -createwallet
./bc.exe -createwallet
./bc.exe -createwallet
./bc.exe -listAllAddresses
