package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func printHelp() {
	fmt.Println("\nAvailable commands:")
	fmt.Println("⬆️ upload <filename>        - Upload a document")
	fmt.Println("⬇️ download <id> <filename> - Download a document")
	fmt.Println("📝 list                     - List all documents")
	fmt.Println("🧐 view                     - View an encrypted document")
	fmt.Println("🙏 help                     - Show this help message")
	fmt.Println("👋 exit                     - Exit the program")
}

func main() {
	serverAddr := flag.String("server", "localhost:50051", "The server address in the format host:port")
	flag.Parse()

	docClient, err := NewDocumentClient(*serverAddr)
	if err != nil {
		log.Fatalf("☠️  Failed to connect to server, is it running?  %v", err)
	}
	defer docClient.Close()

	log.Printf("🔌 Connected to server at %s\n", *serverAddr)
	printHelp()

	scan := bufio.NewScanner(os.Stdin)
	fmt.Print("\n> ")

	for scan.Scan() {
		cmd := strings.Fields(scan.Text())
		if len(cmd) == 0 {
			fmt.Print("> ")
			continue
		}

		switch cmd[0] {
		case "upload":
			if len(cmd) != 2 {
				fmt.Println("Usage: upload <path/to/file>")
				fmt.Print("> ")
				continue
			}
			filename := cmd[1]
			content, err := os.ReadFile(filename)
			if err != nil {
				fmt.Printf("❌ Failed to read file: %v\n", err)
				fmt.Print("> ")
				continue
			}

			docId, err := docClient.UploadDocument(filename, content)
			if err != nil {
				fmt.Printf("❌ Failed to upload document: %v\n", err)
				fmt.Print("> ")
				continue
			}
			fmt.Printf("✅ Document uploaded! ID: %s\n", docId)

		case "download":
			if len(cmd) != 3 {
				fmt.Println("Usage: download <id> <filename to save to>")
				fmt.Print("> ")
				continue
			}
			docId := cmd[1]
			filename := cmd[2]

			doc, err := docClient.DownloadDocument(docId)
			if err != nil {
				fmt.Printf("❌ Failed to download document: %v\n", err)
				fmt.Print("> ")
				continue
			}

			if err := os.WriteFile(filename, doc.Content, 0644); err != nil {
				fmt.Printf("❌ Failed to write file: %v\n", err)
				fmt.Print("> ")
				continue
			}
			fmt.Printf("✅ Document downloaded to %s\n", filename)

		case "list":
			docs, err := docClient.ListDocuments()
			if err != nil {
				fmt.Printf("❌ Failed to list documents: %v\n", err)
				fmt.Print("> ")
				continue
			}

			if len(docs) == 0 {
				fmt.Println("📭 No documents found!")
			} else {
				fmt.Println("\n📚 Available documents:")
				for _, doc := range docs {
					if doc != nil {
						fmt.Printf("ID: %s | Name: %s | Created: %s\n",
							doc.Id, doc.Name, doc.CreatedAt)
					}
				}
				fmt.Println()
			}

		case "view":
			if len(cmd) != 2 {
				fmt.Println("Usage: view <id>")
				fmt.Print("> ")
				continue
			}
			docId := cmd[1]

			doc, err := docClient.ViewDocument(docId)

			if err != nil {
				fmt.Printf("❌ Failed to view document: %v\n", err)
				fmt.Print("> ")
				continue
			}
			fmt.Printf("\n📄 Document Details:\n")
			fmt.Printf("ID: %s\n", doc.Id)
			fmt.Printf("Name: %s\n", doc.Name)
			fmt.Printf("Created: %s\n", doc.CreatedAt)
			fmt.Printf("Updated: %s\n", doc.UpdatedAt)
			fmt.Printf("🔒 Encrypted Content:\n")
			fmt.Printf("Hex: %X\n", doc.Content[:64])

		case "help":
			printHelp()

		case "exit":
			fmt.Println("👋 See ya!")
			return

		default:
			fmt.Printf("🚫 Unknown command: %s\nType 'help' for available commands\n", cmd[0])
		}

		fmt.Print("> ")
	}

	if err := scan.Err(); err != nil {
		fmt.Printf("❌ Error reading input: %v\n", err)
	}
}
