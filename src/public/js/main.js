// Modern Linux Forum JavaScript

document.addEventListener("DOMContentLoaded", function () {
    console.log("DOM fully loaded");

    // For testing - always show the popup regardless of localStorage
    // Remove this line after confirming it works
    showAuthPopup();

    // Initialize the application
    initApp();

    // Add event listeners
    addEventListeners();
});

// Initialize the application
function initApp() {
    console.log("LinuxHub forum initialized");

    // Show authentication popup for first time visitors
    // Comment out the localStorage check for now to ensure popup always shows
    // if (!localStorage.getItem('visited')) {
    showAuthPopup();
    //     localStorage.setItem('visited', 'true');
    // }
}

// Add event listeners to interactive elements
function addEventListeners() {
    // Post creation
    const postInput = document.querySelector('.input-area input');
    if (postInput) {
        postInput.addEventListener('focus', function () {
            this.parentElement.parentElement.classList.add('active');
        });
    }

    // Post button
    const postButton = document.querySelector('.btn.primary');
    if (postButton) {
        postButton.addEventListener('click', function () {
            const input = document.querySelector('.input-area input');
            if (input.value.trim() !== '') {
                createNewPost(input.value);
                input.value = '';
            } else {
                showError("Please enter some content for your post.");
            }
        });
    }

    // Like buttons
    const likeButtons = document.querySelectorAll('.btn-action:first-child');
    likeButtons.forEach(button => {
        button.addEventListener('click', function () {
            toggleLike(this);
        });
    });

    // Menu buttons
    const menuButtons = document.querySelectorAll('.btn-menu');
    menuButtons.forEach(button => {
        button.addEventListener('click', function (e) {
            e.stopPropagation();
            showPostMenu(this);
        });
    });

    // Document click to close menus
    document.addEventListener('click', function () {
        const activeMenu = document.querySelector('.post-menu.active');
        if (activeMenu) {
            activeMenu.remove();
        }
    });
}

// Create a new post from user input
function createNewPost(content) {
    const feed = document.querySelector('.feed');
    const postCreation = document.querySelector('.post-creation');

    // Create post element
    const post = document.createElement('article');
    post.className = 'post card';

    // Get user avatar
    const userAvatar = document.querySelector('.input-area img').src;

    // Current time
    const now = new Date();
    const timeString = 'Just now';

    // Create post HTML
    post.innerHTML = `
        <div class="post-header">
            <img src="${userAvatar}" alt="User" class="user-avatar">
            <div class="post-info">
                <h3>You</h3>
                <span class="post-meta">User • ${timeString}</span>
            </div>
            <button class="btn-menu"><i class="fas fa-ellipsis-h"></i></button>
        </div>
        <div class="post-content">
            <p>${content}</p>
        </div>
        <div class="post-actions">
            <button class="btn-action"><i class="far fa-thumbs-up"></i> Like</button>
            <button class="btn-action"><i class="far fa-comment"></i> Comment</button>
            <button class="btn-action"><i class="far fa-share-square"></i> Share</button>
        </div>
        <div class="post-stats">
            <span>0 likes</span>
            <span>0 comments</span>
        </div>
    `;

    // Add event listeners to new post
    post.querySelector('.btn-action:first-child').addEventListener('click', function () {
        toggleLike(this);
    });

    post.querySelector('.btn-menu').addEventListener('click', function (e) {
        e.stopPropagation();
        showPostMenu(this);
    });

    // Add post to feed
    feed.insertBefore(post, postCreation.nextSibling);

    // Show confirmation
    showNotification("Post created successfully!");

    // Add animation
    post.style.opacity = '0';
    post.style.transform = 'translateY(20px)';
    post.style.transition = 'opacity 0.3s, transform 0.3s';

    setTimeout(() => {
        post.style.opacity = '1';
        post.style.transform = 'translateY(0)';
    }, 10);
}

// Toggle like on posts
function toggleLike(button) {
    const liked = button.classList.contains('liked');
    const postStats = button.closest('.post').querySelector('.post-stats span:first-child');
    const currentLikes = parseInt(postStats.textContent);

    if (liked) {
        // Unlike
        button.classList.remove('liked');
        button.innerHTML = '<i class="far fa-thumbs-up"></i> Like';
        postStats.textContent = `${currentLikes - 1} likes`;
    } else {
        // Like
        button.classList.add('liked');
        button.style.color = 'var(--primary-color)';
        button.innerHTML = '<i class="fas fa-thumbs-up"></i> Like';
        postStats.textContent = `${currentLikes + 1} likes`;

        // Add like animation
        const likeEffect = document.createElement('div');
        likeEffect.className = 'like-effect';
        likeEffect.innerHTML = '<i class="fas fa-thumbs-up"></i>';
        button.appendChild(likeEffect);

        setTimeout(() => {
            if (likeEffect.parentNode) {
                likeEffect.parentNode.removeChild(likeEffect);
            }
        }, 800);
    }
}

// Show post menu (three dots)
function showPostMenu(button) {
    // Remove any existing menus
    const existingMenus = document.querySelectorAll('.post-menu');
    existingMenus.forEach(menu => menu.remove());

    // Create menu
    const menu = document.createElement('div');
    menu.className = 'post-menu active';
    menu.innerHTML = `
        <ul>
            <li><i class="fas fa-bookmark"></i> Save post</li>
            <li><i class="fas fa-bell-slash"></i> Turn off notifications</li>
            <li><i class="fas fa-flag"></i> Report post</li>
        </ul>
    `;

    // Position menu
    const rect = button.getBoundingClientRect();
    menu.style.position = 'absolute';
    menu.style.top = `${rect.bottom + 5}px`;
    menu.style.right = `${window.innerWidth - rect.right}px`;
    menu.style.backgroundColor = 'var(--card-bg)';
    menu.style.boxShadow = '0 2px 10px rgba(0, 0, 0, 0.2)';
    menu.style.borderRadius = '8px';
    menu.style.padding = '8px 0';
    menu.style.zIndex = '1000';

    // Style menu items
    menu.querySelectorAll('ul')[0].style.listStyle = 'none';
    menu.querySelectorAll('li').forEach(item => {
        item.style.padding = '8px 16px';
        item.style.cursor = 'pointer';
        item.style.fontSize = '14px';
        item.style.display = 'flex';
        item.style.alignItems = 'center';
        item.style.color = 'var(--text-primary)';

        item.querySelector('i').style.marginRight = '8px';
        item.style.color = 'var(--text-primary)';

        // Add hover effect
        item.addEventListener('mouseenter', function () {
            this.style.backgroundColor = 'var(--button-hover)';
        });
        item.addEventListener('mouseleave', function () {
            this.style.backgroundColor = 'transparent';
        });

        // Add click handlers
        item.addEventListener('click', function (e) {
            e.stopPropagation();
            menu.remove();
            showNotification("Action applied successfully!");
        });
    });

    document.body.appendChild(menu);

    // Event delegation for menu items
    menu.addEventListener('click', function (e) {
        e.stopPropagation();
    });
}

// Show authentication popup (replaces welcome message)
function showAuthPopup() {
    console.log("Attempting to show auth popup");

    // First check if popup already exists to prevent duplicates
    if (document.querySelector('.auth-popup')) {
        console.log("Auth popup already exists, not creating another one");
        return;
    }

    const authPopup = document.createElement('div');
    authPopup.className = 'auth-popup';
    authPopup.innerHTML = `
        <div class="auth-content card">
            <div class="auth-header">
                <i class="fab fa-linux"></i>
                <h2>Welcome to LinuxHub!</h2>
                <p>The community for Linux enthusiasts and problem-solvers.</p>
            </div>
            
            <div class="auth-tabs">
                <button class="tab-btn active" data-tab="login">Login</button>
                <button class="tab-btn" data-tab="signup">Sign Up</button>
            </div>
            
            <div class="auth-form-container">
                <!-- Login Form -->
                <form class="auth-form login-form active">
                    <div class="form-group">
                        <label for="login-email">Email</label>
                        <input type="email" id="login-email" placeholder="Enter your email">
                    </div>
                    
                    <div class="form-group">
                        <label for="login-password">Password</label>
                        <input type="password" id="login-password" placeholder="Enter your password">
                    </div>
                    
                    <div class="form-footer">
                        <a href="#" class="forgot-password">Forgot password?</a>
                    </div>
                    
                    <button type="submit" class="btn primary btn-block">Login</button>
                </form>
                
                <!-- Signup Form -->
                <form class="auth-form signup-form">
                    <div class="form-group">
                        <label for="signup-email">Email</label>
                        <input type="email" id="signup-email" placeholder="Enter your email">
                    </div>
                    
                    <div class="form-group">
                        <label for="signup-password">Password</label>
                        <input type="password" id="signup-password" placeholder="Create a password">
                    </div>
                    
                    <div class="form-group">
                        <label for="signup-confirm">Confirm Password</label>
                        <input type="password" id="signup-confirm" placeholder="Confirm your password">
                    </div>
                    
                    <button type="submit" class="btn primary btn-block">Create Account</button>
                </form>
            </div>
            
            <div class="guest-option">
                <p>Not ready to sign up?</p>
                <button class="btn btn-text btn-block">Continue as Guest</button>
            </div>
        </div>
    `;

    // Add inline styles temporarily to diagnose issues
    authPopup.style.position = 'fixed';
    authPopup.style.top = '0';
    authPopup.style.left = '0';
    authPopup.style.right = '0';
    authPopup.style.bottom = '0';
    authPopup.style.backgroundColor = 'rgba(0, 0, 0, 0.8)';
    authPopup.style.display = 'flex';
    authPopup.style.alignItems = 'center';
    authPopup.style.justifyContent = 'center';
    authPopup.style.zIndex = '9999';

    // Add to DOM
    document.body.appendChild(authPopup);
    console.log("Auth popup added to DOM");

    // Tab switching functionality
    const tabBtns = authPopup.querySelectorAll('.tab-btn');
    tabBtns.forEach(btn => {
        btn.addEventListener('click', function () {
            // Update tab buttons
            tabBtns.forEach(b => {
                b.classList.remove('active');
            });

            this.classList.add('active');

            // Show corresponding form
            const tabName = this.getAttribute('data-tab');
            const forms = authPopup.querySelectorAll('.auth-form');

            forms.forEach(form => {
                form.classList.remove('active');
            });

            const activeForm = authPopup.querySelector(`.${tabName}-form`);
            activeForm.classList.add('active');
        });
    });

    // Form submission handling (for prototype only)
    const loginForm = authPopup.querySelector('.login-form');
    loginForm.addEventListener('submit', function (e) {
        e.preventDefault();
        const email = document.getElementById('login-email').value;

        if (email) {
            document.body.removeChild(authPopup);
            showNotification(`Logged in successfully as ${email}`);
        } else {
            showError('Please fill in all fields');
        }
    });

    const signupForm = authPopup.querySelector('.signup-form');
    signupForm.addEventListener('submit', function (e) {
        e.preventDefault();
        const email = document.getElementById('signup-email').value;
        const password = document.getElementById('signup-password').value;
        const confirm = document.getElementById('signup-confirm').value;

        if (email && password && confirm) {
            if (password !== confirm) {
                showError('Passwords do not match');
                return;
            }

            document.body.removeChild(authPopup);
            showNotification(`Account created successfully for ${email}`);
        } else {
            showError('Please fill in all fields');
        }
    });

    // Guest option
    const guestButton = authPopup.querySelector('.guest-option .btn');
    guestButton.addEventListener('click', function () {
        document.body.removeChild(authPopup);
        showNotification('Welcome! You are browsing as a guest');
    });
}

// Show notification
function showNotification(message, type = 'success') {
    const notification = document.createElement('div');
    notification.className = 'notification ' + type;
    notification.textContent = message;

    // Style notification
    notification.style.position = 'fixed';
    notification.style.bottom = '20px';
    notification.style.right = '20px';
    notification.style.backgroundColor = type === 'success' ? 'var(--online-color)' : '#e41e3f';
    notification.style.color = 'white';
    notification.style.padding = '10px 20px';
    notification.style.borderRadius = '4px';
    notification.style.boxShadow = '0 2px 10px rgba(0, 0, 0, 0.2)';
    notification.style.zIndex = '9999';

    // Add to DOM
    document.body.appendChild(notification);

    // Remove after 3 seconds
    setTimeout(() => {
        notification.style.opacity = '0';
        notification.style.transform = 'translateY(10px)';
        notification.style.transition = 'opacity 0.3s, transform 0.3s';

        setTimeout(() => {
            if (notification.parentNode) {
                notification.parentNode.removeChild(notification);
            }
        }, 300);
    }, 3000);
}

// Show error message
function showError(message) {
    showNotification(message, 'error');
}

// Add this helper function to force showing the popup
function forceShowAuthPopup() {
    showAuthPopup();
}
