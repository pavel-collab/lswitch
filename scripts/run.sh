#!/bin/bash

active_window=$(xdotool getactivewindow)

xdotool key --window $active_window --clearmodifiers ctrl+c
sleep 0.3

../bin/lswitch -clipboard -quiet

sleep 0.5
xdotool windowactivate $active_window
sleep 0.1

xdotool key --window $active_window --clearmodifiers ctrl+v
