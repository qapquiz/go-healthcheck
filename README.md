# go-healthcheck
This program will check the websites that still healthy or not when given the CSV file. It will report to the hiring line website after finished a health check. You will need to login to the Line platform to send the report.
## Usage
```
go-healthechk test.csv
```
This runs a health check for every website in CSV file. If it is not timeout then will consider as a success.

Output:
```
Perform website checking...
Done!

Checked websites: 75
Successful websites: 54
Failure websites: 21
Total times to finished checking website: 30006 ms

Login with Line to send a report to hiring line. Please allow and login in your browser.
Login successfully
Sending report...
Send report successfully!
```
