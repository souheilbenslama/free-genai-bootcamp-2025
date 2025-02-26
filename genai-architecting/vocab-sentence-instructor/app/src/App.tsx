import React, { useState, useEffect } from 'react';
import styled from 'styled-components';
import TopicSearch from './components/TopicSearch';
import VocabList from './components/VocabList';
import { generateVocabulary } from './services/geminiService';
import { VocabItem } from './types/vocab';
import './App.css';

const AppContainer = styled.div`
  max-width: 800px;
  margin: 0 auto;
  padding: 2rem 1rem;
`;

const Header = styled.header`
  text-align: center;
  margin-bottom: 2rem;
`;

const Title = styled.h1`
  color: #1a73e8;
  margin-bottom: 0.5rem;
`;

const Subtitle = styled.p`
  color: #666;
  font-size: 1.1rem;
  margin-bottom: 2rem;
`;

const ErrorMessage = styled.div`
  background-color: #fdeded;
  color: #b71c1c;
  padding: 1rem;
  border-radius: 8px;
  margin: 1rem 0;
`;

const DebugSection = styled.div`
  margin-top: 2rem;
  padding: 1rem;
  background-color: #f5f5f5;
  border-radius: 8px;
  font-family: monospace;
  font-size: 0.9rem;
`;

function App() {
  const [vocabs, setVocabs] = useState<VocabItem[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [currentTopic, setCurrentTopic] = useState<string | null>(null);
  const [debugInfo, setDebugInfo] = useState<string>('');
  const [showDebug, setShowDebug] = useState<boolean>(false);

  useEffect(() => {
    // Check if API key is set
    const apiKey = process.env.REACT_APP_GEMINI_API_KEY;
    setDebugInfo(prev => prev + `\nAPI Key status: ${apiKey ? 'Set' : 'Not set'}`);
    
    if (!apiKey) {
      setError('API key is not set. Please check your .env file and restart the application.');
    }
  }, []);

  const handleSearch = async (topic: string) => {
    setLoading(true);
    setError(null);
    setCurrentTopic(topic);
    setDebugInfo(`Starting search for topic: ${topic}\n`);
    
    try {
      setDebugInfo(prev => prev + `\nCalling Gemini API...`);
      const response = await generateVocabulary(topic);
      setDebugInfo(prev => prev + `\nReceived response from Gemini API`);
      setVocabs(response.vocabs);
    } catch (err: any) {
      console.error('Error fetching vocabulary:', err);
      setDebugInfo(prev => prev + `\nError: ${err.message || 'Unknown error'}`);
      setError(err.message || 'Failed to generate vocabulary. Please check your API key and try again.');
      setVocabs([]);
    } finally {
      setLoading(false);
    }
  };

  const toggleDebug = () => {
    setShowDebug(!showDebug);
  };

  return (
    <AppContainer>
      <Header>
        <Title>German Vocabulary Instructor</Title>
        <Subtitle>
          Enter a topic in English to get related vocabulary in German
        </Subtitle>
      </Header>

      <TopicSearch onSearch={handleSearch} loading={loading} />
      
      {error && (
        <ErrorMessage>
          <strong>Error:</strong> {error}
          <div style={{ marginTop: '0.5rem' }}>
            <button onClick={toggleDebug}>
              {showDebug ? 'Hide Debug Info' : 'Show Debug Info'}
            </button>
          </div>
        </ErrorMessage>
      )}
      
      {showDebug && (
        <DebugSection>
          <h3>Debug Information</h3>
          <p>Environment: {process.env.NODE_ENV}</p>
          <p>API Key Status: {process.env.REACT_APP_GEMINI_API_KEY ? 'Set (not showing for security)' : 'Not set'}</p>
          <pre>{debugInfo}</pre>
        </DebugSection>
      )}
      
      {currentTopic && !loading && !error && vocabs.length > 0 && (
        <div style={{ marginTop: '1.5rem', textAlign: 'center' }}>
          <h2>Vocabulary for: <span style={{ color: '#1a73e8' }}>{currentTopic}</span></h2>
        </div>
      )}
      
      <VocabList vocabs={vocabs} loading={loading} />
    </AppContainer>
  );
}

export default App;
