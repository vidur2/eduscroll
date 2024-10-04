package embeddingfunction_test

import (
	"fmt"
	"testing"

	"github.com/joho/godotenv"
	embeddingfunction "github.com/vidur2/vectorMicroservices/pkg/embeddingFunction"
)

func TestEmbeddingFunctionApi(t *testing.T) {
	godotenv.Load("../../.env")
	embed := embeddingfunction.EmbedFunction{}
	embeds, err := embed.CreateEmbedding([]string{`The Enchanted Library

		In the small town of Eldoria, hidden between rolling hills and ancient forests, there stood an old, mysterious library that was said to hold the knowledge of the ages. The townspeople spoke in hushed whispers about its enchanted books and the guardian spirit that protected its secrets.
		
		Nina, a curious young girl with a passion for stories, had always been fascinated by the tales surrounding the library. One day, fueled by a burning desire for adventure, she decided to venture into the heart of Eldoria's enchanted woods to find the hidden entrance to the legendary library.
		
		As she navigated through the dense foliage, the air grew thick with magic. The leaves whispered ancient incantations, guiding Nina towards her destiny. Finally, she stumbled upon a moss-covered stone archway, partially concealed by ivy. With a sense of trepidation and excitement, she stepped through.
		
		The library revealed itself in all its grandeur. Shelves upon shelves stretched into the impossibly high ceiling, holding books that glowed with an ethereal light. A soft hum echoed through the vast chamber, and the air was thick with the scent of ancient parchment.
		
		Eventually, Nina realized that her time in the enchanted library had come to an end. With a heart full of gratitude, she bid farewell to the guardian spirit and stepped back through the archway into Eldoria.
		
		As she emerged into the sunlight, clutching the "The Symphony of Time" to her chest, she knew that the stories within were not just tales to be read but seeds of inspiration that would bloom throughout her life. The enchantment of the library had become a part of her, and she carried its magic into the world beyond.
		
		And so, the legend of the enchanted library in Eldoria continued to weave its spell, drawing in those with a thirst for knowledge and a love for the extraordinary. Each visitor added a new chapter to its ever-expanding story, ensuring that the magic endured through the ages.
	`})

	if err != nil {
		fmt.Println(err)
		t.Fatal(err)
	}

	fmt.Println(embeds)
}
