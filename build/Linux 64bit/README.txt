NOTE: Don't move the alert.mp3 file from the executable file's directory. Otherwise script will stop when it tries to play alert

Linux Guide:

    Don't execute dogeElonBot directly or you won't be able to close the program without terminal magic :|
    
    To execute the bot. Open terminal in this directory and just type ./dogeElonBot
    Press Ctrl+C in the same terminal to close it.

    In case you executed the file directly, without the terminal, here are the instructions to close it:
        - Type "ps aux | grep dogeElonBot" in the terminal
        - Note the process id (pid) of dogeElonBot
        - Enter "kill <pid>" in the terminal with the <pid> being the pid you noticed above