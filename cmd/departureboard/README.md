# SimSig Departure Boards

This is a Proof of Concept for a Station Arrivals/Departure Board generated from SimSig WTT files and Interface Gateway Messages.

## Changelog

### v0.2.1 28.01.22

- Do not show train as cancelled at future stops if it is reported as passing a timing point but there was no delay report

### v0.2.0 23.01.22

- Command Line Interface replaced with Graphical Interface
- Convert TIPLOC Names on Location List and if necessary in Origin/Destination

#### Minor Changes
- Add Application Icon
- Replace "Unknown" with "On Time" and "Cancelled"
- Decrease refresh time from 60s to 15s
- Update Clock on Location List Screen
- Add Return Link to Board
- Improve Location List Styling
- Seperate Origin and Destination

### v0.0.1 14.01.22

Initial Release

## Instructions

1. Launch SimSig and start a Simulation, being sure to Enable the "Interface Gateway" on the "Primary" port.
2. Launch departureboard.exe
3. Enter the Simsig Username and Password if running a licensed simulation
4. Click "Choose Timetable" and select the timetable that is in use in SimSig
5. Click "Connect" and wait whilst a connection is established to SimSig
4. Click the "Open Departure Board" link

### Customisation

The template files for the location list and depature board are loaded from `tmpl\index.tmpl` and `tmpl\board.tmpl` respectively.
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

 ## Attribution

Bulletin board icons created by Freepik - [Flaticon](https://www.flaticon.com/free-icons/flight-information)

Extracts from the [Network Rail Reference Data](https://www.networkrail.co.uk/who-we-are/transparency-and-ethics/transparency/open-data-feeds/) are included under the [Open Government License](https://www.nationalarchives.gov.uk/doc/open-government-licence/version/3/)