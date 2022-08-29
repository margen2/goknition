# Goknition:

Goknition is a Go application that uses the Amazon Rekognition API to organize images. It searches for a given collection of faces inside each image and returns a list of matches between them.
  
# Instalation:  

## AWS config:  

You must have AWS configured on your computer. Refer to this link for a detailed explanation: [Getting started](https://aws.amazon.com/getting-started/).  
The Region was set to "us-east-1" as it is one of the cheapest regions that supports the Rekognition API.  
  
**Please be aware that using this application may result in fees to your AWS account.**

## MySQL:
You must have MySQL configured on your computer. You also need to have a .env file with the following values:  

DB_USER=yourUser  
DB_PASSWORD=yourPW    
DB_NAME=goknition  
API_PORT=8080  

Enter the following command on the MySQL monitor:  
`
\. your\path\to\github.com\margen2\goknition\sql\schema.sql
`  

## Cloning the repository:  
`$ git clone https://github.com/margen2/goknition`  
`$ cd goknition`  
`$ go run main.go`  
  
Now open your browser on http://localhost:8080/
 
# How to use:

## File structure:  

The files are expected to be distributed as follows:    

<pre>      
│
├──IDs                   
│  ├──ID1
│  │  └──file1.JPG   
│  ├──ID2
│  │  └──file2.JPG
│  └──ID3
│     └──file2.JPG
│
└── DATA
    ├──folder1
    │    └──subfolder
    │        └──file1.JPG     
    ├──folder2
    │  │──file2.JPG
    │  │──file3.JPG
    │  └──file4.JPG
    └──folder3
       └──file5.JPG 
</pre>  
  

The IDs folder **Must** be organized as shown here. Each subfolder's name will be used to query for a particular face. The Data folder doesn't have those same restrictions as any file inside it will be saved. 

## Creating a collection: 
A Rekognition face collection will hold all the faces you want to identify in the DATA folder. To create one, select the "Create collection" option on the menu and insert the path to the IDs folder and a name for the collection. 

## Uploading images:
To search for faces inside an image, select the "Upload images" option and insert the path to the DATA folder with the name of the collection you want to use.

## Results:
After uploading your images, you can search by a specific face by selecting the "Make a query" option, or you can save all of the matches inside another folder using the "Save matches" option.