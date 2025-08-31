#!/bin/sh


export ANDROID_HOME="{$HOME}/Android/Sdk/"
export ANDROID_NDK_HOME="${HOME}/Android/Sdk/ndk/29.0.13599879"

cd cmd/mepm-gui

 fyne package --os android/arm64  --app-id com.melias.mepm --name "mepm"

 #  adb install -r mepm.apk

 # adb shell monkey -p com.melias.mepm 1
