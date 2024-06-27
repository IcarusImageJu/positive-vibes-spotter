# Positive-Vibes-Spotter

## Project Description:

Positive-Vibes-Spotter is a script designed for my cyberdeck that captures photos of different places and generates positive descriptions about them. The aim is to create uplifting messages that help people feel loved and appreciated when they read them on the cyberdeck. Whether it’s a bustling cityscape or a quiet corner of nature, Positive-Vibes-Spotter finds the beauty in every scene and shares it with a smile.



https://github.com/IcarusImageJu/positive-vibes-spotter/assets/10758538/e0a1d509-8433-48d7-a4c3-f81f5c97252c



## Features:

  - Captures photos automatically at regular intervals.
  - Analyzes the captured images to generate positive and encouraging descriptions.
  - Displays the messages on the cyberdeck to spread positivity and good vibes.
  - Support french only for now

## Usage:

  1. Copy the `spot.sh` file to your cyberdeck.
  2. create a cron job to run the script at regular intervals.
  3. Let the script run and enjoy the positive messages generated from the places around you.

## CRON JOB EXAMPLE:

To automate the execution of the script, you need to set up a cron job on your cyberdeck and ensure that the OPENAI_API_KEY environment variable is correctly set. Below is an example of how to do this:

  1. Open the crontab editor by typing `crontab -e` in your terminal.
  2. Add the following lines to schedule the script:

```
@reboot OPENAI_API_KEY="votre_clé_api" ~/spot.sh  # Runs the script once at reboot with the API key
0 7-22 * * * OPENAI_API_KEY="votre_clé_api" ~/spot.sh  # Runs the script every hour from 7 AM to 10 PM with the API key
```

## OpenAI API Key:

The script requires an OpenAI API key to generate the positive descriptions. You can obtain your API key by signing up on the OpenAI platform and generating a key from your account settings. Ensure to keep your API key secure and do not share it publicly.

## Tech Stack:

  - BASH
  - Debian

## Future Enhancements:

  - Add Ollama support for local use only
  - Customizable message templates.
  - Upgrade UI Rendering
  - Add multilingual support

## Contributing:

Contributions are welcome! Please feel free to submit issues or pull requests.
