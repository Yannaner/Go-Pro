// API Base URL
const API_URL = 'http://localhost:8080';

// MediaPipe Hands
let hands;
let camera;
const videoElement = document.getElementById('video');
const canvasElement = document.getElementById('canvas');
const canvasCtx = canvasElement.getContext('2d');

// Buttons
const startCameraBtn = document.getElementById('startCamera');
const analyzeHandBtn = document.getElementById('analyzeHand');
const sendChatBtn = document.getElementById('sendChat');
const chatInput = document.getElementById('chatInput');
const chatContainer = document.getElementById('chatContainer');
const submitCardsBtn = document.getElementById('submitCards');
const cardInputSection = document.getElementById('cardInputSection');

// Card inputs
const card1Rank = document.getElementById('card1Rank');
const card1Suit = document.getElementById('card1Suit');
const card2Rank = document.getElementById('card2Rank');
const card2Suit = document.getElementById('card2Suit');

// Results
const handDetectedSpan = document.querySelector('#handDetected span');
const winRateSpan = document.querySelector('#winRate span');

// Hand detection state
let handDetected = false;
let currentHand = null;

// Initialize MediaPipe Hands
function initMediaPipe() {
    hands = new Hands({
        locateFile: (file) => {
            return `https://cdn.jsdelivr.net/npm/@mediapipe/hands/${file}`;
        }
    });

    hands.setOptions({
        maxNumHands: 2,
        modelComplexity: 1,
        minDetectionConfidence: 0.5,
        minTrackingConfidence: 0.5
    });

    hands.onResults(onResults);
}

// Process MediaPipe results
function onResults(results) {
    // Clear canvas
    canvasCtx.save();
    canvasCtx.clearRect(0, 0, canvasElement.width, canvasElement.height);
    canvasCtx.drawImage(results.image, 0, 0, canvasElement.width, canvasElement.height);

    if (results.multiHandLandmarks) {
        handDetected = true;
        handDetectedSpan.textContent = 'Yes';
        handDetectedSpan.style.color = '#4CAF50';

        for (const landmarks of results.multiHandLandmarks) {
            drawConnectors(canvasCtx, landmarks, HAND_CONNECTIONS, {
                color: '#00FF00',
                lineWidth: 5
            });
            drawLandmarks(canvasCtx, landmarks, {
                color: '#FF0000',
                lineWidth: 2
            });
        }
    } else {
        handDetected = false;
        handDetectedSpan.textContent = 'No';
        handDetectedSpan.style.color = '#ff4444';
    }

    canvasCtx.restore();
}

// Start camera
startCameraBtn.addEventListener('click', async () => {
    try {
        const stream = await navigator.mediaDevices.getUserMedia({
            video: { width: 640, height: 480 }
        });
        
        videoElement.srcObject = stream;
        
        // Set canvas size to match video
        canvasElement.width = 640;
        canvasElement.height = 480;
        
        initMediaPipe();
        
        camera = new Camera(videoElement, {
            onFrame: async () => {
                await hands.send({ image: videoElement });
            },
            width: 640,
            height: 480
        });
        
        camera.start();
        
        startCameraBtn.disabled = true;
        startCameraBtn.textContent = 'Camera Running';
        analyzeHandBtn.disabled = false;
        
    } catch (error) {
        alert('Error accessing camera: ' + error.message);
    }
});

// Analyze hand
analyzeHandBtn.addEventListener('click', async () => {
    if (!handDetected) {
        alert('No hand detected! Please show your poker cards to the camera.');
        return;
    }

    // Show card input section
    cardInputSection.style.display = 'block';
    analyzeHandBtn.textContent = 'Waiting for cards...';
    analyzeHandBtn.disabled = true;
});

// Submit cards for analysis
submitCardsBtn.addEventListener('click', async () => {
    // Validate card inputs
    if (!card1Rank.value || !card1Suit.value || !card2Rank.value || !card2Suit.value) {
        alert('Please select both cards (rank and suit)');
        return;
    }

    try {
        const handData = {
            cards: [
                { rank: card1Rank.value, suit: card1Suit.value },
                { rank: card2Rank.value, suit: card2Suit.value }
            ]
        };

        const response = await fetch(`${API_URL}/api/analyze`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(handData)
        });

        if (!response.ok) {
            throw new Error('Analysis failed');
        }

        const result = await response.json();
        currentHand = result;
        
        winRateSpan.textContent = `${(result.winrate * 100).toFixed(1)}%`;
        winRateSpan.style.color = result.winrate > 0.5 ? '#4CAF50' : '#ff9800';
        
        // Reset button
        analyzeHandBtn.textContent = 'Analyze Hand';
        analyzeHandBtn.disabled = false;
        
        // Automatically get AI insights after analysis
        await getAIInsights(result);
        
    } catch (error) {
        alert('Error analyzing hand: ' + error.message);
        analyzeHandBtn.textContent = 'Analyze Hand';
        analyzeHandBtn.disabled = false;
    }
});

// Send chat message
async function sendMessage() {
    const message = chatInput.value.trim();
    if (!message) return;

    // Add user message to chat
    addMessage(message, 'user');
    chatInput.value = '';

    try {
        const response = await fetch(`${API_URL}/api/chat`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                message: message,
                hand: currentHand
            })
        });

        if (!response.ok) {
            throw new Error('Chat request failed');
        }

        const result = await response.json();
        addMessage(result.response, 'bot');
        
    } catch (error) {
        addMessage('Sorry, I encountered an error: ' + error.message, 'bot');
    }
}

// Add message to chat
function addMessage(text, sender) {
    const messageDiv = document.createElement('div');
    messageDiv.className = `message ${sender}-message`;
    messageDiv.innerHTML = `<strong>${sender === 'user' ? 'You' : 'AI Coach'}:</strong> ${text}`;
    chatContainer.appendChild(messageDiv);
    chatContainer.scrollTop = chatContainer.scrollHeight;
}

// Get AI insights automatically after hand analysis
async function getAIInsights(hand) {
    if (!hand || !hand.cards || hand.cards.length === 0) {
        return;
    }

    // Build a descriptive message about the hand
    const cardsDescription = hand.cards.map(c => `${c.rank} of ${c.suit}`).join(', ');
    const winRatePercent = (hand.winrate * 100).toFixed(1);
    
    const autoMessage = `I have ${cardsDescription} with a ${winRatePercent}% win rate. What's your analysis and strategy recommendation?`;
    
    // Show the auto-generated question in chat
    addMessage(autoMessage, 'user');

    try {
        const response = await fetch(`${API_URL}/api/chat`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                message: autoMessage,
                hand: hand
            })
        });

        if (!response.ok) {
            throw new Error('Chat request failed');
        }

        const result = await response.json();
        addMessage(result.response, 'bot');
        
    } catch (error) {
        addMessage('Sorry, I encountered an error getting insights: ' + error.message, 'bot');
    }
}

// Event listeners
sendChatBtn.addEventListener('click', sendMessage);
chatInput.addEventListener('keypress', (e) => {
    if (e.key === 'Enter') {
        sendMessage();
    }
});

// Disable analyze button initially
analyzeHandBtn.disabled = true;
