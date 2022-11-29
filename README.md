# Goknition:

Goknition is a Go desktop application that uses the Amazon Rekognition API to organize images. It searches for a given collection of faces inside each image and returns a list of matches between them.
  
# Instalation:  

## AWS config:  

You must have AWS configured on your computer. Refer to this link for a detailed explanation: [Getting started](https://aws.amazon.com/getting-started/).  
The Region was set to "us-east-1" as it is one of the cheapest regions that supports the Rekognition API.  
  
**Please be aware that using this application may result in fees to your AWS account.**

## MySQL:
Make sure you have [MySQL](https://dev.mysql.com/downloads/installer/) installed

## Wails:
Wails provides a way to create desktop application with go, follow the [Getting started](https://wails.io/docs/gettingstarted/installation/) guide to install it.
 
 
# How to use:

First, clone the repository:

`$ git clone https://github.com/margen2/goknition`  
`$ cd goknition`  
`$ wails build`  

Open the gokniton.exe file created after running wails build. The first time you run the program, it will create a config.json file, open it and update the db_password and db_user fields with your credentials. Now open the gokniton.exe file again.


## Creating a collection: 
A Rekognition face collection will hold the faces you want to use for identification. To create one, select the "Collections" option on the menu and click on "create". Then, click on "add faces" to add the faces folder and click on "select" to set it as the active collection. 

## Searching Images:
To scan the image folder and identify face matches, select the "Images" option and click on "Data folder" to select the folder containing the images you want to scan. 

## Results:
After uploading your images, select "Save matches" and click on the "save matches" button next to the collection, then choose the folder where you want to save the matches.

## File structure:  

The images are expected to be distributed as follows:    

<pre>      
│
├──IDs                   
│  ├──ID1
│  │  │──file1.JPG
│  │  └──file2.JPG
│  │   
│  ├──ID2
│  │  │──file1.JPG
│  │  │──file2.JPG
│  │  └──file3.JPG
│  └──ID3
│     └──file1.JPG
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
  
The IDs folder **Must** be organized as shown here. Each subfolder's name will be used to query for a particular face. The Data folder doesn't have those same restrictions as any file inside it will be saved. Check [here](https://docs.aws.amazon.com/rekognition/latest/dg/recommendations-facial-input-images.html) for best facial recognition practices.
