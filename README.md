usage:
`go build -o build/app && ./build/app`

takes `input.json`, which is array of jsons

```json
{
    "type": 0-5,
    "name":"something",
    "result":"something",
    "params":{
        "file_name":"foo.txt", 
        "new_file_name":"bar.txt",//optional, needed for new filename
        "message":"Inferno", // optional, needed for Writing to file
        "compare_date": ""
    }
}

```
Types: 

- 0 - create file
- 1 - change filename 
- 2 - delete file 
- 3 - get file creation date and time 
- 4 - write to file
- 5 - check if date and time more than  

    date format: "YYYY-MM-DD HH:MM:SS -0700",-0700 for timezone

    ```json 
    [
        {
            "type":5,
            //...
        },
        {
            //action 1
        },
        {
            //action 2 
        },
    ]
    ```
    Если дата и время, указанные больше чем у файла, указанного в параметрах, то происходит выполнение action 1 и пропускается выполнение action 2, если же дата и время меньше или равно, то происходит пропуск action 1 и выполняется action 2.
