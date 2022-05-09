Scaffolding for a tool that reads ids from a csv file a sends concurrent requests to a base url using them. Then it writes the responses into another csv file. 

The idea is to use some of the concepts behind goroutines and channels. 

You can pass the number of goroutines (workers) as an argument. It defaults to 1.

Example usage (should maintain the ids file structure)
```bash
$ go build
$ ./datagrabber -w 8
```
