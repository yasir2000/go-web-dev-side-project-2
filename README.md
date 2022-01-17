# go-web-dev-side-project-2

**Command-line tools**

The application we've built so far is ready to take the world by storm but not before we give it a home on the Internet. We need to pick a valid, catchy, and available domain name, which we can point to the server running our Go code. Instead of sitting in front of our favorite domain name provider for hours on end trying different names, we are going to develop a few command-line tools that will help us find the right one. As we do so, we will see how the Go standard library allows us to interface with the terminal and other executing applications; we'll also explore some patterns and practices to build command-line programs.

In this section, you will learn:



* How to build complete command-line applications with as little as a single code file
* How to ensure that the tools we build can be composed with other tools using standard streams
* How to interact with a simple third-party JSON RESTful API
* How to utilize the standard in and out pipes in Go code
* How to read from a streaming source, one line at a time
* How to build a WHOIS client to look up domain information
* How to store and use sensitive or deployment-specific information in environment variables

<h5>**Pipe design for command-line tools**</h5>


We are going to build a series of command-line tools that use the standard streams (`stdin` and `stdout`) to communicate with the user and with other tools. Each tool will take an input line by line via the standard input pipe, process it in some way, and then print the output line by line to the standard out pipe for the next tool or user.

By default, the standard input is connected to the user's keyboard, and the standard output is printed to the terminal from where the command was run; however, both can be redirected using **redirection metacharacters**. It's possible to throw the output away by redirecting it to `NUL` on Windows or `/dev/null` on Unix machines, or redirecting it to a file that will cause the output to be saved to a disk. Alternatively, you can pipe (using the `|` pipe character) the output of one program to the input of another; it is this feature that we will make use of in order to connect our various tools together. For example, you could pipe the output from one program to the input of another program in a terminal using this code:



1. `echo -n "Hello" | md5`

The output of the `echo` command will be the string `Hello` (without the quotes), which is then **piped** to the `md5` command; this command will in turn calculate the MD5 hash of `Hello`:



1. `8b1a9953c4611296a827abf8c47804d7`

Our tools will work with lines of strings where each line (separated by a linefeed character) represents one string. When run without any pipe redirection, we will be able to interact directly with the programs using the default in and out, which will be useful when testing and debugging our code.

<h5>**Five simple programs**</h5>


In this section, we will build five small programs that we will combine at the end. The key features of the programs are as follows:



* **Sprinkle**: This program will add some web-friendly sprinkle words to increase the chances of finding the available domain names.
* **Domainify**: This program will ensure words are acceptable for a domain name by removing unacceptable characters. Once this is done, it will replace spaces with hyphens and add an appropriate top-level domain (such as `.com` and `.net`) to the end.
* **Coolify**: This program will change a boring old normal word to Web 2.0 by fiddling around with vowels.
* **Synonyms**: This pro will use a third-party API to find synonyms.
* **Available**: This gram will use a third-party API to find synonyms. Available: This program will check to see whether the domain is available or not using an appropriate **WHOIS** server.

Five programs might seem like a lot for one section, but don't forget how small entire programs can be in Go.

<h5>**Sprinkle**</h5>


Our first program augments the incoming words with some sugar terms in order to improve the odds of finding names that are available. Many companies use this approach to keep the core messaging consistent while being able to afford the `.com` domain. For example, if we pass in the word `chat`, it might pass out `chatapp`; alternatively, if we pass in `talk`, we may get back `talk time`.

Go's `math/rand` package allows us to break away from the predictability of computers. It gives our program the appearance of intelligence by introducing elements of chance into its decision making.

To make our Sprinkle program work, we will:



* Define an array of transformations, using a special constant to indicate where the original word will appear
* Use the `bufio` package to scan the input from `stdin` and `fmt.Println` in order to write the output to `stdout`
* Use the `math/rand` package to randomly select a transformation to apply

**Tip**:** **All our programs will reside in the `$GOPATH/src` directory. For example, if your GOPATH is `~/Work/projects/go`, you would create your program folders in the `~/Work/projects/go/src` folder.

In the `$GOPATH/src` directory, create a new folder called `sprinkle` and add a `main.go` file containing the following code:


``` bash
`package main`
import (`
 "bufio"`
 "fmt"`
 "math/rand"`
 "os"`
 "strings"`
 "time"`
)
const otherWord = "*"`
var transforms = []string{ 
otherWord,
otherWord + "app",
otherWord + "site",
otherWord + "time",
"get" + otherWord,
"go" + otherWord,
"lets " + otherWord,
 otherWord + "hq",
}
func main() {
 rand.Seed(time.Now().UTC().UnixNano())
 s := bufio.NewScanner(os.Stdin)
 for s.Scan() {
   t := transforms[rand.Intn(len(transforms))]
  fmt.Println(strings.Replace(t, otherWord, s.Text(), -1))
 }
}
``` 

From now on, it is assumed that you will sort out the appropriate `import`statements yourself.

The preceding code represents our complete Sprinkle program. It defines three things: a constant, a variable, and the obligatory `main` function, which serves as the entry point to Sprinkle. The `otherWord` constant string is a helpful token that allows us to specify where the original word should occur in each of our possible transformations. It lets us write code, such as `otherWord+"extra"`, which makes it clear that in this particular case, we want to add the word "extra" to the end of the original word.

The possible transformations are stored in the `transforms` variable that we declare as a slice of strings. In the preceding code, we defined a few different transformations, such as adding `app` to the end of a word or `lets` before it. Feel free to add some more; the more creative, the better.

In the `main` function, the first thing we do is use the current time as a random seed. Computers can't actually generate random numbers, but changing the seed number of random algorithms gives the illusion that it can. We use the current time in nanoseconds because it's different each time the program is run (provided the system clock isn't being reset before each run). If we skip this step, the numbers generated by the `math/rand` package would be deterministic; they'd be the same every time we run the program.

We then create a `bufio.Scanner` object (by calling `bufio.NewScanner`) and tell it to read the input from `os.Stdin`, which represents the standard input stream. This will be a common pattern in our five programs since we are always going to read from the standard **in** and write to the standard **out**.

**Tip**:** **The `bufio.Scanner` object actually takes `io.Reader` as its input source, so there is a wide range of types that we could use here. If you were writing unit tests for this code, you could specify your own `io.Reader` for the scanner to read from, removing the need for you to worry about simulating the standard input stream.

As the default case, the scanner allows us to read blocks of bytes separated by defined delimiters, such as carriage return and linefeed characters. We can specify our own split function for the scanner or use one of the options built in the standard library. For example, there is `bufio.ScanWords`, which scans individual words by breaking on whitespace rather than linefeeds. Since our design specifies that each line must contain a word (or a short phrase), the default line-by-line setting is ideal.

A call to the `Scan` method tells the scanner to read the next block of bytes (the next line) from the input, and then it returns a `bool` value indicating whether it found anything or not. This is how we are able to use it as the condition for the `for` loop. While there is content to work on, `Scan` returns `true` and the body of the `for` loop is executed; when `Scan` reaches the end of the input, it returns `false`, and the loop is broken. The bytes that are selected are stored in the `Bytes` method of the scanner, and the handy `Text` method that we use converts the `[]byte` slice into a string for us.

Inside the `for` loop (so for each line of input), we use `rand.Intn` to select a random item from the transforms slice and use `strings.Replace` to insert the original word where the `otherWord` string appears. Finally, we use `fmt.Println` to print the output to the default standard output stream.

**Note**: The `math/rand` package provides insecure random numbers. If you want to write code that utilizes random numbers for security purposes, you must use the `crypto/rand` package instead.

Let's build our program and play with it:



1. `go build -o sprinkle`
2. `./sprinkle`

Once the program starts running, it will use the default behavior to read the user input from the terminal. It uses the default behavior because we haven't piped in any content or specified a source for it to read from. Type `chat` and hit **Return**. The scanner in our code notices the linefeed character at the end of the word and runs the code that transforms it, outputting the result. For example, if you type `chat` a few times, you would see the following output:



1. `chat`
2. `go chat`
3. `chat`
4. `lets chat`
5. `chat`
6. `chat app`

Sprinkle never exits (meaning the `Scan` method never returns `false` to break the loop) because the terminal is still running; in normal execution, the in pipe will be closed by whatever program is generating the input. To stop the program, hit **Ctrl** + **C**.

Before we move on, let's try to run Sprinkle, specifying a different input source. We are going to use the `echo` command to generate some content and pipe it to our Sprinkle program using the pipe character:



1. `echo "chat" | ./sprinkle`

The program will randomly transform the word, print it out, and exit since the `echo` command generates only one line of input before terminating and closing the pipe.

We have successfully completed our first program, which has a very simple but useful function, as we will see.

**Tip**: As an extra assignment, rather than hardcoding the `transformations` array as we have done, see whether you can externalize it via flags or store them in a text file or database.

<h5>**Domainify**</h5>


Some of the words that output from Sprinkle contain spaces and perhaps other characters that are not allowed in domains. So we are going to write a program called Domainify; it converts a line of text into an acceptable domain segment and adds an appropriate **Top-level Domain** (**TLD**) to the end. Alongside the `sprinkle` folder, create a new one called `domainify` and add the `main.go` file with the following code:


``` bash
package main
var tlds = []string{"com", "net"}
const allowedChars = "abcdefghijklmnopqrstuvwxyz0123456789_-"
func main() {
 rand.Seed(time.Now().UTC().UnixNano())
 s := bufio.NewScanner(os.Stdin)
for s.Scan() {
  text := strings.ToLower(s.Text())
   var newText []rune
   for _, r := range text {
     if unicode.IsSpace(r) {
       r = '-'
     }
     if !strings.ContainsRune(allowedChars, r) {
       continue
     }
     newText = append(newText, r)
   }
     fmt.Println(string(newText) + "." +        
               tlds[rand.Intn(len(tlds))])
  }
}
```

You will notice a few similarities between Domainify and the Sprinkle program: we set the random seed using `rand.Seed`, generate a `NewScanner`  method wrapping the `os.Stdin` reader, and scan each line until there is no more input.

We then convert the text to lowercase and build up a new slice of `rune` types called `newText`. The `rune` types consist of only characters that appear in the `allowedChars`  string, which `strings.ContainsRune` lets us know. If `rune` is a space that we determine by calling `unicode.IsSpace`, we replace it with a hyphen, which is an acceptable practice in domain names.

**Note:** Ranging over a string returns the index of each character and a rune type, which is a numerical value (specifically, int32) representing the character itself. For more information about runes, characters, and strings, refer to [http://blog.golang.org/strings](http://blog.golang.org/strings).

Finally, we convert `newText` from a `[]rune` slice into a string and add either `.com` or `.net` at the end, before printing it out using `fmt.Println`.

Let's build and run Domainify:



1. `go build -o domainify`
2. `./domainify`

Type in some of these options to see how `domainify` reacts:



* `Monkey`
* `Hello Domainify`
* `"What's up?"`
* `One (two) three!`

You can see that, for example, `One (two) three!` might yield `one-two-three.com`.

We are now going to compose Sprinkle and Domainify to see them work together. In your terminal, navigate to the parent folder (probably `$GOPATH/src`) of sprinkle and domainify and run the following command:



1. `./sprinkle/sprinkle | ./domainify/domainify`

Here, we ran the `sprinkle` program and piped the output to the `domainify`program. By default, `sprinkle` uses the terminal as the input and `domanify`outputs to the terminal. Try typing in `chat` a few times again and notice the output is similar to what Sprinkle was outputting previously, except now they are acceptable for domain names. It is this piping between programs that allows us to compose command-line tools together.

**Tip**: Only supporting .com and .net top-level domains is fairly limiting. As an additional assignment, see whether you can accept a list of TLDs via a command-line flag.

<h5>**Coolify**</h5>


Often, domain names for common words, such as `chat`, are already taken, and a common solution is to play around with the vowels in the words. For example, we might remove `a` and make it `cht` (which is actually less likely to be available) or add `a` to produce `chaat`. While this clearly has no actual effect on coolness, it has become a popular, albeit slightly dated, way to secure domain names that still sound like the original word.

Our third program, Coolify, will allow us to play with the vowels of words that come in via the input and write modified versions to the output.

Create a new folder called `coolify` alongside `sprinkle` and `domainify`, and create the `main.go` code file with the following code:


``` bash
package main
const (
 duplicateVowel bool   = true
 removeVowel    bool   = false
) 
func randBool() bool {
 return rand.Intn(2) == 0
}
func main() {
 rand.Seed(time.Now().UTC().UnixNano())
 s := bufio.NewScanner(os.Stdin)
 for s.Scan() {
   word := []byte(s.Text())
   if randBool() {
     var vI int = -1
     for i, char := range word {
       switch char {
       case 'a', 'e', 'i', 'o', 'u', 'A', 'E', 'I', 'O', 'U':
   if randBool() {
          vI = i
         }
      }
     }
    if vI >= 0 {
       switch randBool() {
       case duplicateVowel:
         word = append(word[:vI+1], word[vI:]...)
       case removeVowel:
         word = append(word[:vI], word[vI+1:]...)
       }
    }
   }
   fmt.Println(string(word))
 }
 }
```

While the preceding Coolify code looks very similar to the code of Sprinkle and Domainify, it is slightly more complicated. At the very top of the code, we declare two constants, `duplicateVowel` and `removeVowel`, that help make the Coolify code more readable. The `switch` statement decides whether we duplicate or remove a vowel. Also, using these constants, we are able to express our intent very clearly, rather than use just `true` or `false`.

We then define the `randBool` helper function that just randomly returns either `true` or `false`. This is done by asking the `rand` package to generate a random number and confirming whether that number comes out as zero. It will be either `0` or `1`, so there's a fifty-fifty chance of it being `true`.

The `main` function of Coolify starts the same way as that of Sprinkle and Domainify setting the `rand.Seed` method and creating a scanner of the standard input stream before executing the loop body for each line of input. We call `randBool` first to decide whether we are even going to mutate a word or not, so Coolify will only affect half the words passed through it.

We then iterate over each rune in the string and look for a vowel. If our `randBool` method returns `true`, we keep the index of the vowel character in the `vI` variable. If not, we keep looking through the string for another vowel, which allows us to randomly select a vowel from the words rather than always modify the same one.

Once we have selected a vowel, we use `randBool` again to randomly decide what action to take.

**Note**: This is where the helpful constants come in; consider the following alternative switch statement:



1. `switch randBool() {`
2. ` case true:`
3. `   word = append(word[:vI+1], word[vI:]...)`
4. ` case false:`
5. `   word = append(word[:vI], word[vI+1:]...) }`

In the preceding code snippet, it's difficult to tell what is going on because `true ` and `false ` don't express any context. On the other hand, using `duplicateVowel ` and `removeVowel ` tells anyone reading the code what we mean by the result of `randBool` .

The three dots following slices cause each item to pass as a separate argument to the `append` function. This is an idiomatic way of appending one slice to another. Inside the `switch` case, we do some slice manipulation to either duplicate the vowel or remove it altogether. We are slicing our `[]byte` slice again and using the `append` function to build a new one made up of sections of the original word. The following diagram shows which sections of the string we access in our code:



<p id="gdcalert1" ><span style="color: red; font-weight: bold">>>>>>  gd2md-html alert: error handling inline image </span><br>(<a href="#">Back to top</a>)(<a href="#gdcalert2">Next alert</a>)<br><span style="color: red; font-weight: bold">>>>>> </span></p>



If we take the value `blueprints` as an example word and assume that our code has selected the first `e` character as the vowel (so that `vI` is `3`), the following table will illustrate what each new slice of the word will represent:


<table>
  <tr>
   <td><strong>Code</strong>
   </td>
   <td><strong>Value</strong>
   </td>
   <td><strong>Description</strong>
   </td>
  </tr>
  <tr>
   <td><code>word[:vI+1]</code>
   </td>
   <td>blue
   </td>
   <td>This describes the slice from the beginning of the word until the selected vowel. The <code>+1</code> is required because the value following the colon does not include the specified index; rather, it slices up to that value.
   </td>
  </tr>
  <tr>
   <td><code>word[vI:]</code>
   </td>
   <td>eprints
   </td>
   <td>This describes the slice starting from and including the selected vowel to the end of the slice.
   </td>
  </tr>
  <tr>
   <td><code>word[:vI]</code>
   </td>
   <td>blu
   </td>
   <td>This describes the slice from the beginning of the word up to, but not including, the selected vowel.
   </td>
  </tr>
  <tr>
   <td><code>word[vI+1:]</code>
   </td>
   <td>prints
   </td>
   <td>This describes the slice from the item following the selected vowel to the end of the slice.
   </td>
  </tr>
</table>


After we modify the word, we print it out using `fmt.Println`.

Let's build Coolify and play with it to see what it can do:



1. `go build -o coolify`
2. `./coolify`

When Coolify is running, try typing `blueprints` to see what sort of modifications it comes up with:



1. `blueprnts`
2. `bleprints`
3. `bluepriints`
4. `blueprnts`
5. `blueprints`
6. `bluprints`

Let's see how Coolify plays with Sprinkle and Domainify by adding their names to our pipe chain. In the terminal, navigate back (using the `cd`command) to the parent folder and run the following commands:



1. `./coolify/coolify | ./sprinkle/sprinkle | ./domainify/domainify`

We will first spice up a word with extra pieces and make it cooler by tweaking the vowels before finally transforming it into a valid domain name. Play around by typing in a few words and seeing what suggestions our code makes.

**Tip**: Coolify only works on vowels; as an additional exercise, see whether you can make the code operate on every character it encounters just to see what happens.

<h5>**Synonyms**</h5>


So far, our programs have only modified words, but to really bring our solution to life, we need to be able to integrate a third-party API that provides word synonyms. This allows us to suggest different domain names while retaining the original meaning. Unlike Sprinkle and Domainify, Synonyms will write out more than one response for each word given to it. Our architecture of piping programs together means this won't be much of a problem; in fact, we do not even have to worry about it since each of the three programs is capable of reading multiple lines from the input source.

Big Huge Thesaurus, [http://bighugelabs.com/](http://bighugelabs.com/), has a very clean and simple API that allows us to make a single HTTP `GET` request to look up synonyms.

Before you can use Big Huge Thesaurus, you'll need an API key, which you can get by signing up to the service at [http://words.bighugelabs.com/](http://words.bighugelabs.com/).

<h5>**Using environment variables for configuration**</h5>


Your API key is a sensitive piece of configuration information that you don't want to share with others. We could store it as `const` in our code. However, this would mean we will not be able to share our code without sharing our key (not good, especially if you love open source projects). Additionally, perhaps more importantly, you will have to recompile your entire project if the key expires or if you want to use a different one (you don't want to get into such a situation).

A better solution is using an environment variable to store the key, as this will allow you to easily change it if you need to. You could also have different keys for different deployments; perhaps you could have one key for development or testing and another for production. This way, you can set a specific key for a particular execution of code so you can easily switch between keys without having to change your system-level settings. Also, different operating systems deal with environment variables in similar ways, so they are a perfect choice if you are writing cross-platform code.

Create a new environment variable called `BHT_APIKEY` and set your API key as its value.

**Note**: For machines running a bash shell, you can modify your `~/.bashrc` file or similar to include `export` commands, such as the following:



1. `export BHT_APIKEY=abc123def456ghi789jkl`

On Windows machines, you can navigate to the properties of your computer and look for Environment Variables in the Advanced section.

<h5>**Consuming a web API**</h5>


Making a request for in a web browser shows us what the structure of JSON response data looks like when finding synonyms for the word `love`:

``` bash

{
 "noun":{
   "syn":[
     "passion",
     "beloved",
     "dear"
   ]
 },
 "verb":{
  "syn":[
     "love",
     "roll in the hay",
     "make out"
   ],
   "ant":[
     "hate"
   ]
 }
}
```

A real API will return a lot more actual words than what is printed here, but the structure is the important thing. It represents an object, where the keys describe the types of word (verbs, nouns, and so on). Also, values are objects that contain arrays of strings keyed on `syn` or `ant` (for the synonym and antonym, respectively); it is the synonyms we are interested in.

To turn this JSON string data into something we can use in our code, we must decode it into structures of our own using the capabilities found in the `encoding/json` package. Because we're writing something that could be useful outside the scope of our project, we will consume the API in a reusable package rather than directly in our program code. Create a new folder called `thesaurus` alongside your other program folders (in `$GOPATH/src`) and insert the following code into a new `bighuge.go` file:

```bash

package thesaurus
import (
 "encoding/json"
 "errors"
 "net/http"
)
type BigHuge struct {
 APIKey string
}
type synonyms struct {
 Noun *words `json:"noun"`
 Verb *words `json:"verb"`
}
type words struct {
 Syn []string `json:"syn"`
}
func (b *BigHuge) Synonyms(term string) ([]string, error) {
 var syns []string
 response, err := http.Get("http://words.bighugelabs.com/api/2/"  +
  b.APIKey + "/" + term + "/json")
 if err != nil {
   return syns, errors.New("bighuge: Failed when looking for  synonyms   
    for "" + term + """ + err.Error())
 }
 var data synonyms
 defer response.Body.Close()
 if err := json.NewDecoder(response.Body).Decode(&data); err !=  nil {
   return syns, err
 }
 if data.Noun != nil {
   syns = append(syns, data.Noun.Syn...)
 }
 if data.Verb != nil {
   syns = append(syns, data.Verb.Syn...)
 }
 return syns, nil
}
```

In the preceding code, the `BigHuge` type we define houses the necessary API key and provides the `Synonyms` method that will be responsible for doing the work of accessing the endpoint, parsing the response, and returning the results. The most interesting parts of this code are the `synonyms` and `words`structures. They describe the JSON response format in Go terms, namely an object containing noun and verb objects, which in turn contain a slice of strings in a variable called `Syn`. The tags (strings in backticks following each field definition) tell the `encoding/json` package which fields to map to which variables; this is required since we have given them different names.

**Tip**: Typically in JSON, keys have lowercase names, but we have to use capitalized names in our structures so that the `encoding/json` package would also know that the fields exist. If we don't, the package would simply ignore the fields. However, the types themselves (`synonyms` and `words`) do not need to be exported.

The `Synonyms` method takes a `term` argument and uses `http.Get` to make a web request to the API endpoint in which the URL contains not only the API key value, but also the `term` value itself. If the web request fails for some reason, we will make a call to `log.Fatalln`, which will write the error to the standard error stream and exit the program with a non-zero exit code (actually an exit code of `1`). This indicates that an error has occurred.

If the web request is successful, we pass the response body (another `io.Reader`) to the `json.NewDecoder` method and ask it to decode the bytes into the `data` variable that is of our `synonyms` type. We defer the closing of the response body in order to keep the memory clean before using Go's built-in `append` function to concatenate both `noun` and `verb` synonyms to the `syns`slice that we then return.

Although we have implemented the `BigHuge` thesaurus, it isn't the only option out there, and we can express this by adding a `Thesaurus` interface to our package. In the `thesaurus` folder, create a new file called `thesaurus.go`and add the following interface definition to the file:



1. `package thesaurus`
2. `type Thesaurus interface {`
3. ` Synonyms(term string) ([]string, error)`
4. `}`

This simple interface just describes a method that takes a `term` string and returns either a slice of strings containing the synonyms or an error (if something goes wrong). Our `BigHuge` structure already implements this interface, but now, other users could add interchangeable implementations for other services, such as [http://www.dictionary.com/](http://www.dictionary.com/) or the Merriam-Webster online service.

Next, we are going to use this new package in a program. Change the directory in the terminal back up a level to `$GOPATH/src`, create a new folder called `synonyms`, and insert the following code into a new `main.go` file you will place in this folder:


``` bash
func main() {
 apiKey := os.Getenv("BHT_APIKEY")
 thesaurus := &thesaurus.BigHuge{APIKey: apiKey}
 s := bufio.NewScanner(os.Stdin)
 for s.Scan() {
   word := s.Text()
   syns, err := thesaurus.Synonyms(word)
   if err != nil {
     log.Fatalln("Failed when looking for synonyms for  "+word+", err)
   }
   if len(syns) == 0 {
     log.Fatalln("Couldn't find any synonyms for " + word +  ")
   }
   for _, syn := range syns {
     fmt.Println(syn)
   }
 }
}
```

Now when you manage your imports again, you will have written a complete program that is capable of looking up synonyms of words by integrating the Big Huge Thesaurus API.

In the preceding code, the first thing our `main` function does is that it gets the `BHT_APIKEY` environment variable value via the `os.Getenv` call. To protect your code, you might consider double-checking it to ensure the value is properly set; if not, report the error. For now, we will assume that everything is configured properly.

Next, the preceding code starts to look a little familiar since it scans each line of input again from `os.Stdin` and calls the `Synonyms` method to get a list of the replacement words.

Let's build a program and see what kind of synonyms the API comes back with when we input the word `chat`:



1. `go build -o synonyms`
2. `./synonyms`
3. `chat`
4. `confab`
5. `confabulation`
6. `schmooze`
7. `New World chat`
8. `Old World chat`
9. `conversation`
10. `thrush`
11. `wood warbler`
12. `chew the fat`
13. `shoot the breeze`
14. `chitchat`
15. `chatter`

The results you get will most likely differ from what we have listed here since we're hitting a live API. However, the important thing is that when we provide a word or term as an input to the program, it returns a list of synonyms as the output, one per line.

<h5>**Getting domain suggestions**</h5>


By composing the four programs we have built so far in this section, we already have a useful tool for suggesting domain names. All we have to do now is to run the programs while piping the output to the input in an appropriate way. In a terminal, navigate to the parent folder and run the following single line:



1. `./synonyms/synonyms | ./sprinkle/sprinkle | ./coolify/coolify |  ./domainify`

Because the `synonyms` program is first in our list, it will receive the input from the terminal (whatever the user decides to type in). Similarly, because `domainify` is last in the chain, it will print its output to the terminal for the user to see. Along the way, the lines of words will be piped through other programs, giving each of them a chance to do their magic.

Type in a few words to see some domain suggestions; for example, when you type `chat` and hit **Return**, you may see the following:



1. `getcnfab.com`
2. `confabulationtim.com`
3. `getschmoozee.net`
4. `schmosee.com`
5. `neew-world-chatsite.net`
6. `oold-world-chatsite.com`
7. `conversatin.net`
8. `new-world-warblersit.com`
9. `gothrush.net`
10. `lets-wood-wrbler.com`
11. `chw-the-fat.com`

The number of suggestions you get will actually depend on the number of synonyms. This is because it is the only program that generates more lines of output than what we input.

We still haven't solved our biggest problem: the fact that we have no idea whether the suggested domain names are actually available or not. So we still have to sit and type each one of them into a website. In the next section, we will address this issue.

<h5>Available</h5>


Our final program, Available, will connect to a WHOIS server to ask for details about the domains passed to it of course, if no details are returned, we can safely assume that the domain is available for purchase. Unfortunately, the WHOIS specification (see [http://tools.ietf.org/html/rfc3912](http://tools.ietf.org/html/rfc3912)) is very small and contains no information about how a WHOIS server should reply when you ask for details about a domain. This means programmatically parsing the response becomes a messy endeavor. To address this issue for now, we will integrate with only a single WHOIS server, which we can be sure will have `No match` somewhere in the response when it has no records for the domain.

**Note**: A more robust solution is to have a WHOIS interface with a well-defined structure for the details and perhaps an error message for cases when the domain doesn't exist with different implementations for different WHOIS servers. As you can imagine, it's quite a project; it is perfect for an open source effort.

Create a new folder called `available` alongside others and add a `main.go` file to it containing the following function code:

``` bash

func exists(domain string) (bool, error) {
const whoisServer string = "com.whois-servers.net"
 conn, err := net.Dial("tcp", whoisServer+":43")
 if err != nil {
   return false, err
 }
 defer conn.Close()
 conn.Write([]byte(domain + "rn"))
 scanner := bufio.NewScanner(conn)
 for scanner.Scan() {
   if strings.Contains(strings.ToLower(scanner.Text()), "no match") {
     return false, nil
   }
 }
 return true, nil
}
```

The `exists` function implements what little there is in the WHOIS specification by opening a connection to port `43` on the specified `whoisServer`instance with a call to `net.Dial`. We then defer the closing of the connection, which means that no matter how the function exits (successful, with an error, or even a panic), `Close()` will still be called on the `conn` connection. Once the connection is open, we simply write the domain followed by `rn` (the carriage return and linefeed characters). This is all that the specification tells us, so we are on our own from now on.

Essentially, we are looking for some mention of "no match" in the response, and this is how we will decide whether a domain exists or not (`exists` in this case is actually just asking the WHOIS server whether it has a record for the domain we specified). We use our favorite `bufio.Scanner` method to help us iterate over the lines in the response. Passing the connection to `NewScanner`works because `net.Conn` is actually an `io.Reader` too. We use `strings.ToLower`so we don't have to worry about case sensitivity and `strings.Contains` to check whether any one of the lines contains the `no match` text. If it does, we return `false` (since the domain doesn't exist); otherwise, we return `true`.

The `com.whois-servers.net` WHOIS service supports domain names for `.com`and `.net`, which is why the Domainify program only adds these types of domains. If you had used a server that had WHOIS information for a wider selection of domains, you could have added support for additional TLDs.

Let's add a `main` function that uses our `exists` function to check whether the incoming domains are available or not. The check mark and cross mark symbols in the following code are optional if your terminal doesn't support them you are free to substitute them with simple `Yes` and `No` strings.

Add the following code to `main.go`:


``` bash

var marks = map[bool]string{true: "✔", false: "✖"}
func main() {
s := bufio.NewScanner(os.Stdin)
for s.Scan() {
domain := s.Text()
fmt.Print(domain, " ")
exist, err := exists(domain)
if err != nil {
log.Fatalln(err)
}
fmt.Println(marks[!exist])
time.Sleep(1 * time.Second)
}
}
```

**Note**: We can use the check and cross characters in our code happily because all Go code files are UTF-8 compliant the best way to actually get these characters is to search the Web for them and use the copy and paste option to bring them into our code. Otherwise, there are platform-dependent ways to get such special characters.

In the preceding code for the `main` function, we simply iterate over each line coming in via `os.Stdin`. This process helps us print out the domain with `fmt.Print` (but not `fmt.Println`, as we do not want the linefeed yet), call our `exists` function to check whether the domain exists or not, and print out the result with `fmt.Println` (because we do want a linefeed at the end).

Finally, we use `time.Sleep` to tell the process to do nothing for a second in order to make sure we take it easy on the WHOIS server.

**Tip**: Most WHOIS servers will be limited in various ways in order to prevent you from taking up too much in terms of resources. So, slowing things down is a sensible way to make sure we don't make the remote servers angry.

Consider what this also means for unit tests. If a unit test were actually making real requests to a remote WHOIS server, every time your tests run, you will be clocking up statistics against your IP address. A much better approach would be to stub the WHOIS server to simulate responses.

The `marks` map at the top is a nice way to map the `bool` response from `exists` to human-readable text, allowing us to just print the response in a single line using `fmt.Println(marks[!exist])`. We are saying not exist because our program is checking whether the domain is available or not (logically, the opposite of whether it exists in the WHOIS server or not).

After fixing the import statements for the main.go file, we can try out Available to see whether the domain names are available or not by typing the following command:



1. `go build -o available`
2. `./available`

Once Available is running, type in some domain names and see the result appear on the next line:



<p id="gdcalert2" ><span style="color: red; font-weight: bold">>>>>>  gd2md-html alert: error handling inline image </span><br>(<a href="#">Back to top</a>)(<a href="#gdcalert3">Next alert</a>)<br><span style="color: red; font-weight: bold">>>>>> </span></p>



As you can see, for domains that are not available, we get a little cross mark next to them; however, when we make up a domain name using random numbers, we see that it is indeed available.

<h5>**Composing all five programs**</h5>


Now that we have completed all five programs, it's time to put them all together so that we can use our tool to find an available domain name for our chat application. The simplest way to do this is to use the technique we have been using throughout this section: using pipes in a terminal to connect the output and input.

In the terminal, navigate to the parent folder of the five programs and run the following single line of code:



1. `./synonyms/synonyms | ./sprinkle/sprinkle | ./coolify/coolify |  ./domainify/domainify | ./available/available`

Once the programs are running, type in a starting word and see how it generates suggestions before checking their availability.

For example, typing in `chat` might cause the programs to take the following actions:

1. The word `chat` goes into `synonyms`, which results in a series of synonyms:



* confab
* confabulation
* schmooze

2. The synonyms flow into `sprinkle`; here they are augmented with web-friendly prefixes and suffixes, such as the following:



* confabapp
* goconfabulation
* schmooze time

3. These new words flow into `coolify`; here the vowels are potentially tweaked:



* confabaapp
* goconfabulatioon
* schmoooze time

4. The modified words then flow into `domainify`; here they are turned into valid domain names:



* confabaapp.com
* goconfabulatioon.net
* schmooze-time.com

5. Finally, the domain names flow into `available`; here they are checked against the WHOIS server to see whether somebody has already taken the domain or not:



* confabaapp.com

<p id="gdcalert3" ><span style="color: red; font-weight: bold">>>>>>  gd2md-html alert: error handling inline image </span><br>(<a href="#">Back to top</a>)(<a href="#gdcalert4">Next alert</a>)<br><span style="color: red; font-weight: bold">>>>>> </span></p>


* goconfabulatioon.net 

<p id="gdcalert4" ><span style="color: red; font-weight: bold">>>>>>  gd2md-html alert: error handling inline image </span><br>(<a href="#">Back to top</a>)(<a href="#gdcalert5">Next alert</a>)<br><span style="color: red; font-weight: bold">>>>>> </span></p>


* schmooze-time.com 

<p id="gdcalert5" ><span style="color: red; font-weight: bold">>>>>>  gd2md-html alert: error handling inline image </span><br>(<a href="#">Back to top</a>)(<a href="#gdcalert6">Next alert</a>)<br><span style="color: red; font-weight: bold">>>>>> </span></p>



<h5>**One program to rule them all**</h5>


Running our solution by piping programs together is an elegant form of architecture, but it doesn't have a very elegant interface. Specifically, whenever we want to run our solution, we have to type the long, messy line where each program is listed and separated by pipe characters. In this section, we are going to write a Go program that uses the `os/exec` package to run each subprogram while piping the output from one to the input of the next, as per our design.

Create a new folder called `domainfinder` alongside the other five programs and create another new folder called `lib` inside this folder. The `lib` folder is where we will keep builds of our subprograms, but we don't want to copy and paste them every time we make a change. Instead, we will write a script that builds the subprograms and copies the binaries to the `lib` folder for us.

Create a new file called `build.sh` on Unix machines or `build.bat` for Windows and insert into it the following code:



1. `#!/bin/bash`
2. `echo Building domainfinder...`
3. `go build -o domainfinder`
4. `echo Building synonyms...`
5. `cd ../synonyms`
6. `go build -o ../domainfinder/lib/synonyms`
7. `echo Building available...`
8. `cd ../available`
9. `go build -o ../domainfinder/lib/available`
10. `cd ../build`
11. `echo Building sprinkle...`
12. `cd ../sprinkle`
13. `go build -o ../domainfinder/lib/sprinkle`
14. `cd ../build`
15. `echo Building coolify...`
16. `cd ../coolify`
17. `go build -o ../domainfinder/lib/coolify`
18. `cd ../build`
19. `echo Building domainify...`
20. `cd ../domainify`
21. `go build -o ../domainfinder/lib/domainify`
22. `cd ../build`
23. `echo Done.`

The preceding script simply builds all our subprograms (including `domainfinder`, which we are yet to write), telling `go build` to place them in our `lib` folder. Be sure to give execution rights to the new script by doing `chmod +x build.sh` or something similar. Run this script from a terminal and look inside the `lib` folder to ensure that it has indeed placed the binaries for our subprograms.

**Tip**: Don't worry about the `no buildable Go source files` error for now; it's just Go telling us that the `domainfinder` program doesn't have any `.go` files to build.

Create a new file called `main.go` inside `domainfinder` and insert the following code into the file:


```bash
package main
var cmdChain = []*exec.Cmd{
 exec.Command("lib/synonyms"),
 exec.Command("lib/sprinkle"),
 exec.Command("lib/coolify"),
 exec.Command("lib/domainify"),
 exec.Command("lib/available"),
}
func main() {
 cmdChain[0].Stdin = os.Stdin
 cmdChain[len(cmdChain)-1].Stdout = os.Stdout
 for i := 0; i &lt; len(cmdChain)-1; i++ {
   thisCmd := cmdChain[i]
   nextCmd := cmdChain[i+1]
   stdout, err := thisCmd.StdoutPipe()
   if err != nil {
     log.Fatalln(err)
   }
   nextCmd.Stdin = stdout
 }
 for _, cmd := range cmdChain {
   if err := cmd.Start(); err != nil {
     log.Fatalln(err)
  } else {
     defer cmd.Process.Kill()
   }
 }
 for _, cmd := range cmdChain {
   if err := cmd.Wait(); err != nil {
     log.Fatalln(err)
   }
 }
}
```

The `os/exec` package gives us everything we need to work with to run external programs or commands from within Go programs. First, our `cmdChain` slice contains `*exec.Cmd` commands in the order in which we want to join them together.

At the top of the `main` function, we tie the `Stdin` (standard in stream) of the first program with the `os.Stdin` stream of this program and the `Stdout` (standard out stream) of the last program with the `os.Stdout` stream of this program. This means that, like before, we will be taking input through the standard input stream and writing output to the standard output stream.

Our next block of code is where we join the subprograms together by iterating over each item and setting its `Stdin` to the `Stdout` stream of the program before it.

The following table shows each program with a description of where it gets its input from and where its output goes:


<table>
  <tr>
   <td><strong>Program</strong>
   </td>
   <td><strong>Input (Stdin)</strong>
   </td>
   <td><strong>Output (Stdout)</strong>
   </td>
  </tr>
  <tr>
   <td>synonyms
   </td>
   <td>The same <code>Stdin</code> as <code>domainfinder</code>
   </td>
   <td>sprinkle
   </td>
  </tr>
  <tr>
   <td>sprinkle
   </td>
   <td>synonyms
   </td>
   <td>coolify
   </td>
  </tr>
  <tr>
   <td>coolify
   </td>
   <td>sprinkle
   </td>
   <td>domainify
   </td>
  </tr>
  <tr>
   <td>domainify
   </td>
   <td>coolify
   </td>
   <td>available
   </td>
  </tr>
  <tr>
   <td>available
   </td>
   <td>domainify
   </td>
   <td>The same <code>Stdout</code> as <code>domainfinder</code>
   </td>
  </tr>
</table>


We then iterate over each command calling the `Start` method, which runs the program in the background (as opposed to the `Run` method, which will block our code until the subprogram exists which would be no good since we will have to run five programs at the same time). If anything goes wrong, we bail with `log.Fatalln`; however, if the program starts successfully, we defer a call to kill the process. This helps us ensure the subprograms exit when our `main` function exits, which will be when the `domainfinder` program ends.

Once all the programs start running, we iterate over every command again and wait for it to finish. This is to ensure that `domainfinder` doesn't exit early and kill off all the subprograms too soon.

Run the `build.sh` or `build.bat` script again and notice that the `domainfinder` program has the same behavior as we have seen before, with a much more elegant interface.

The following screenshot shows the output from our programs when we type `clouds`; we have found quite a few available domain name options:



<p id="gdcalert6" ><span style="color: red; font-weight: bold">>>>>>  gd2md-html alert: error handling inline image </span><br>(<a href="#">Back to top</a>)(<a href="#gdcalert7">Next alert</a>)<br><span style="color: red; font-weight: bold">>>>>> </span></p>



<h5>**Summary**</h5>


In this section, we learned how five small command-line programs can, when composed together, produce powerful results while remaining modular. We avoided tightly coupling our programs so they could still be useful in their own right. For example, we can use our Available program just to check whether the domain names we manually enter are available or not, or we can use our `synonyms` program just as a command-line thesaurus.

We learned how standard streams could be used to build different flows of these types of programs and how the redirection of standard input and standard output lets us play around with different flows very easily.

We learned how simple it is in Go to consume a JSON RESTful API web service when we wanted to get the synonyms from Big Huge Thesaurus. We also consumed a non-HTTP API when we opened a connection to the WHOIS server and wrote data over raw TCP.

We saw how the `math/rand` package can bring a little variety and unpredictability by allowing us to use pseudo random numbers and decisions in our code, which means that each time we run our program, we will get different results.

Finally, we built our `domainfinder` super program that composes all the subprograms together, giving our solution a simple, clean, and elegant interface.

In the next section, we will take some ideas we have learned so far one step further by exploring how to connect programs using messaging queue technologies allowing them to distributed across many machines to achieve large scale.
