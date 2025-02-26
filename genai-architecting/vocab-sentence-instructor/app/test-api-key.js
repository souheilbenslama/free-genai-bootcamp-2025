// Simple script to test the Gemini API key
const { GoogleGenerativeAI } = require('@google/generative-ai');

// Use the API key directly for testing
const API_KEY = process.env.REACT_APP_GEMINI_API_KEY;

// Initialize the API client
const genAI = new GoogleGenerativeAI(API_KEY);

async function testApiKey() {
  try {
    console.log('Testing Gemini API key...');
    
    // Use gemini-2.0-flash which is confirmed to be working
    const model = genAI.getGenerativeModel({ model: 'gemini-2.0-flash' });
    
    // Simple prompt to test if the API key works
    const prompt = 'Say hello in German';
    
    console.log('Sending request to Gemini API...');
    const result = await model.generateContent(prompt);
    console.log('Received response from Gemini API');
    
    const response = result.response;
    const text = response.text();
    
    console.log('API Response:', text);
    console.log('API key is working correctly!');
  } catch (error) {
    console.error('Error testing API key:', error);
    console.error('API key may not be valid or may not have access to the Gemini API.');
    
    // List available models
    try {
      console.log('\nTrying to list available models...');
      const models = await genAI.listModels();
      console.log('Available models:', models);
    } catch (listError) {
      console.error('Error listing models:', listError);
    }
  }
}

// Run the test
testApiKey();
