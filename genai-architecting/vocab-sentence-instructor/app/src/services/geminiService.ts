import { GoogleGenerativeAI } from '@google/generative-ai';
import { VocabResponse } from '../types/vocab';

const API_KEY = process.env.REACT_APP_GEMINI_API_KEY || '';
const genAI = new GoogleGenerativeAI(API_KEY);

export const generateVocabulary = async (topic: string): Promise<VocabResponse> => {
  try {
    console.log('Using API Key:', API_KEY ? 'API key is set' : 'API key is missing');
    
    if (!API_KEY) {
      throw new Error('API key is missing. Please check your .env file.');
    }
    
    // Use gemini-2.0-flash which is confirmed to be working
    const model = genAI.getGenerativeModel({ model: 'gemini-2.0-flash' });
    
    const prompt = `
      Generate a list of 10 vocabulary words related to the topic "${topic}" in English with their German translations.
      For each word, provide:
      1. The English word
      2. The German translation
      3. A breakdown of the parts of the German word (if applicable)
      
      Format the response as a valid JSON object with the following structure:
      {
        "vocabs": [
          {
            "english": "example word",
            "german": "Beispielwort",
            "parts": ["Beispiel", "wort"]
          },
          ...
        ]
      }
      
      Only return the JSON object, nothing else.
    `;

    console.log('Sending request to Gemini API...');
    const result = await model.generateContent(prompt);
    console.log('Received response from Gemini API');
    const response = result.response;
    const text = response.text();
    console.log('Raw API response:', text);
    
    // Extract the JSON object from the response
    const jsonMatch = text.match(/\{[\s\S]*\}/);
    if (!jsonMatch) {
      throw new Error('Failed to parse JSON from the API response');
    }
    
    try {
      const jsonResponse = JSON.parse(jsonMatch[0]) as VocabResponse;
      return jsonResponse;
    } catch (parseError) {
      console.error('JSON parsing error:', parseError);
      throw new Error(`Failed to parse JSON: ${parseError instanceof Error ? parseError.message : 'Unknown error'}`);
    }
  } catch (error: unknown) {
    console.error('Error generating vocabulary:', error);
    
    if (error instanceof Error) {
      const errorMessage = error.message || '';
      
      if (errorMessage.includes('PERMISSION_DENIED')) {
        throw new Error('API key permission denied. Please check if your API key is valid and has access to the Gemini API.');
      } else if (errorMessage.includes('INVALID_ARGUMENT')) {
        throw new Error('Invalid request to Gemini API. Please try a different topic.');
      } else if (errorMessage.includes('Not Found') && errorMessage.includes('models/')) {
        throw new Error('The specified Gemini model was not found. Please check the model name or try a different model.');
      } else {
        throw error;
      }
    } else {
      throw new Error('An unknown error occurred while generating vocabulary.');
    }
  }
};
