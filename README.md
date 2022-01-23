# Golang SimSig Tools

This is a collection of tools for working with the Signalling Simulation [SimSig](http://www.simsig.co.uk). It is currently in a very basic state and offers no useful API.

## Command departureboard

A tool which uses WTT data and the Interface Gateway to generate a departure board

## Command jsontest

Example of decoding a JSON Message

## Command corpusconvert

A tool which read Network Rail CORPUS Reference Data in JSON Format and writes out a list of TIPLOC/Location Description Pairs in MsgPack format

## Submodule wttxml

- Data structures for Timetable Data in XML format
- Utility functions for accessing the XML data

## Submodule gateway

- Data Structures for unmarshalling Interface Gateway messages.