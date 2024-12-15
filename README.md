# yandex-backup

The program is needed to make backup copies of files to yandex disk

# Variables

| Var                       | Description                                                 |
|---------------------------|:------------------------------------------------------------|
| YANDEX_TOKEN              | Token for access to yandex disk                             |
| BACKUP_DIR                | The directory where you need to look for files for backup   |
| FILE_PREFIX               | How the name of the file to be backed up should begin       |
| YANDEX_DIR                | In which directory on the yandex disk should you put files? |
| NOTIFICATION_CHAT_ID      | сhat_id of the user to whom you want to send notifications  |
| NOTIFICATION_BOT_TOKEN    | Telegram bot token for sending notifications                |
| NOTIFICATION_SUBJECT_LINE | Subject of the notice                                       |

# Error codes

| Code error | Description                                                                                                            |
|------------|:-----------------------------------------------------------------------------------------------------------------------|
| 48045      | Environment variables are not set                                                                                      |
| 4315...    | Incorrect response code when creating a directory on yandex disk. The last three digits are equal to the response code |
| 4316...    | Incorrect response code when get upload url. The last three digits are equal to the response code                      |
| 47046      | Error reading a file to send to disk                                                                                   |
| 11940      | Failed to create a request to send a file                                                                              |
| 53076      | Error when sending a file to disk                                                                                      |
| 73121      | Error converting the form to send notification                                                                         |

