# Business Requirements

- Use generative AI to create personalized language lessons and assistance based on user proficiency and goals.
- Provide interactive exercises and study activities based on the lessons to adapt to the user's learning pace.
- Offer detailed analytics on user progress, strengths, and areas for improvement.
- Implement robust security measures to protect user data.

## Functional Requirements

- Implement AI-driven chatbots for real-time conversation practice in the target language.
- Offer a sentence constructor AI assistance for users to learn creating full, correct sentences in the target language.
- Provide translations and explanations for complex phrases and idioms.
- Track user progress and learning history.
- Consider developing a mobile app for on-the-go learning.

## Assumptions

- We are assuming that the open-source LLMs we choose will be powerful enough to run on hardware with an investment of $10-15k.
- We are just going to hook up a single server in our office to the internet and we should have enough bandwidth to serve 300 students.

## Non-Functional Requirements

- **Performance:** The system should handle up to 300 concurrent users with minimal latency. Response time for AI-driven features should be under 2 seconds.
- **Scalability:** The architecture should support easy scaling to accommodate future growth in user base and features.
- **Availability:** Ensure 99.9% uptime to provide a reliable learning experience.
- **Security:** Implement encryption for data in transit and at rest. Regularly update and patch systems to protect against vulnerabilities.
