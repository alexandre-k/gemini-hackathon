# Gemini Hackathon

## Summary

A project to show a proof of concept of Gemini coupled with a smartphone.

If the permission to access the camera is given, then the user has 2 options:
1. Get an description of what the camera sees
2. Get an English version of a text readable on the camera

Essentially it gives a user a means of using Gemini to get oral interpretations of images.

## How to run

A couple of steps are required to run our web app. Install certificates with
mkcert and run docker compose. If docker is not available on your OS, please
refer to a guide to install it.

```bash
cd keys
mkcert hostname.local
mkcert -install
cd ..
docker compose up
```
