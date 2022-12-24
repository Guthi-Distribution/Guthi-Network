# Guthi Network 
Network part of the **Guthi-Distribution** project done as a final year project of undergraduate level. 

### Requirements
* Language: Go 1.19, [TDM-GCC](https://jmeubank.github.io/tdm-gcc/)
* Build System: CMake, Ninja or Makefile


### Build Instruction 
Create a `config.json` as shown below:
```
{
    "name": "node_name",
    "address": "ip_address" // if you want to use localhost use ""
}
```


* Build and run the executable and use of `-port` parameter is optional and default is port `6969`. 
```
go build
./GuthiNetwork
```