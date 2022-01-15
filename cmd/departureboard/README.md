# SimSig Departure Boards

This is a Proof of Concept for a Station Arrivals/Departure Board generated from SimSig WTT files and Interface Gateway Messages.

## Changelog

### v0.0.1 14.01.22

Initial Release

## Instructions

1. Launch SimSig and start a Simulation, being sure to Enable the "Interface Gateway" on the "Primary" port.
2. Launch departureboard.exe
3. A file dialog will appear, use this to select the timetable that is running in SimSig.
4. A web browser will automatically be launched to access the web interface at http://localhost:8090/

If the software does not function as expected, please consult the error messages printed to the 
console window. If the application exits immediately, it will be necessary to run the application from a "Command Prompt" in order to see these messages.

### Licensed Simulations

In case you are running a licensed simulation, your SimSig username and password must be supplied on the command line.
This can be done either from a "Command Prompt" or by creating a Windows Shortcut.

The command should have the following format: `depatureboard.exe -user <myusername> -pass <mypassword>`

## Command Line Reference
```
  -all
        Do not hide departed and terminated trains
  -help
        Print help text
  -pass string
        SimSig License Password(optional)
  -server string
        Simsig Interface Gateway Address (default "localhost:51515")
  -user string
        SimSig License Username(optional)
  -verbose
        Print received train movement messages
  -wtt string
        Path to Timetable file
```

### Customisation

The template files for the location list and depature board are loaded from tmpl\index.tmpl and board.tmpl respectively.
The visual appearance of the web interface can be customised by modifying these files.

## Author

(c) 2022 Jonathan Pilborough
Contact: jonathan@pilborough.co.uk

## License

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files 
(the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge,
 publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so
 , subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF 
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE 
FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
 WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
