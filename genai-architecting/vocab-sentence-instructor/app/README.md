# German Vocabulary Instructor

A React application that allows users to input a topic in English and get a list of vocabulary words related to that topic in German. Each vocabulary item includes the English word, German translation, and a breakdown of the parts of the German word.

## Features

- Search for vocabulary words by topic
- View German translations for English words
- See the breakdown of compound German words
- Clean and modern user interface

## Technologies Used

- React
- TypeScript
- Styled Components
- Google's Generative AI (Gemini API)

## Getting Started

### Prerequisites

- Node.js (v14 or higher)
- npm or yarn
- Gemini API key

### Installation

1. Clone the repository
2. Navigate to the project directory
3. Install dependencies:

```bash
npm install
```

4. Create a `.env` file in the root directory and add your Gemini API key:

```
REACT_APP_GEMINI_API_KEY=your_gemini_api_key_here
```

### Running the Application

Start the development server:

```bash
npm start
```

The application will be available at [http://localhost:3000](http://localhost:3000).

## How to Use

1. Enter a topic in English in the search bar (e.g., "food", "travel", "technology")
2. Click the "Search" button
3. View the list of vocabulary words related to the topic
4. Each card displays:
   - The English word
   - The German translation
   - A breakdown of the parts of the German word (if applicable)

## Note

You need to obtain a Gemini API key from Google AI Studio to use this application. Replace the placeholder in the `.env` file with your actual API key.

## License

MIT
