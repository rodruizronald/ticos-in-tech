# TicosInTech - Job Board MVP

A modern job board focused on Costa Rica's tech industry, built with React and Tailwind CSS.

## 🚀 Features

### Phase 1 (MVP) - ✅ Completed
- **Landing page** with hero section and search
- **Job listings** with search and filtering
- **Job detail pages** with full job information
- **Responsive design** for mobile, tablet, and desktop
- **Company logos** with fallback initials
- **Save jobs** functionality (localStorage)
- **Modern UI** with Costa Rica-themed design system

### Phase 2 (Coming Soon)
- User registration and authentication
- Email job alerts
- Advanced search features
- Company profiles

## 🛠️ Tech Stack

- **React 18** with functional components and hooks
- **Tailwind CSS** for styling with custom design system
- **Lucide React** for icons
- **Custom router** (lightweight, can upgrade to React Router)
- **Context API** for state management
- **Local Storage** for saved jobs
- **Fetch API** for HTTP requests

## 📋 Prerequisites

- Node.js 16+ and npm
- Running API server at `http://localhost:8080/api/v1/jobs`

## 🚀 Getting Started

### 1. Create the project
```bash
npx create-react-app ticos-in-tech
cd ticos-in-tech
```

### 2. Install dependencies
```bash
npm install lucide-react
npm install -D tailwindcss postcss autoprefixer
npx tailwindcss init -p
```

### 3. Replace default files
Copy all the files from the artifacts into your project structure:

```
src/
├── components/
│   ├── common/
│   │   ├── Header.jsx
│   │   ├── Footer.jsx
│   │   └── LoadingSpinner.jsx
│   ├── job/
│   │   ├── JobCard.jsx
│   │   ├── JobGrid.jsx
│   │   ├── JobDetail.jsx
│   │   └── FilterSidebar.jsx
│   ├── search/
│   │   ├── HeroSection.jsx
│   │   └── SearchBar.jsx
│   └── layout/
│       └── Layout.jsx
├── context/
│   └── AppContext.jsx
├── hooks/
│   ├── useApi.js
│   └── useLocalStorage.js
├── pages/
│   ├── HomePage.jsx
│   ├── JobDetailPage.jsx
│   └── NotFoundPage.jsx
├── services/
│   └── api.js
├── utils/
│   ├── constants.js
│   ├── formatters.js
│   └── router.js
├── styles/
│   └── index.css
├── App.jsx
└── index.js
```

### 4. Update configuration files
- Replace `tailwind.config.js` with the provided configuration
- Replace `public/index.html` with the provided HTML
- Update `package.json` with the provided dependencies

### 5. Start the development server
```bash
npm start
```

The app will open at `http://localhost:3000`

## 🔧 Configuration

### API Configuration
Update the API base URL in `src/utils/constants.js`:
```javascript
export const API_BASE_URL = 'http://localhost:8080/api/v1';
```

### Design System
The app uses a Costa Rica-themed design system defined in:
- `tailwind.config.js` - Colors, fonts, spacing
- `src/styles/index.css` - Component styles and utilities

### Environment Variables
Create a `.env` file for environment-specific settings:
```env
REACT_APP_API_BASE_URL=http://localhost:8080/api/v1
REACT_APP_SITE_NAME=TicosInTech
```

## 📱 Responsive Design

The app is fully responsive with breakpoints:
- **Mobile**: 0-767px (single column, mobile-first)
- **Tablet**: 768-1023px (adapted layouts)
- **Desktop**: 1024px+ (full features, multi-column)

## 🎨 Design System

### Colors
- **Primary**: Costa Rica Blue (#0052CC)
- **Secondary**: Costa Rica Red (#E53E3E)
- **Grays**: 50-900 scale for UI elements

### Typography
- **Font**: Inter (Google Fonts)
- **Scale**: 12px to 48px with defined line heights

### Spacing
- **Base unit**: 4px
- **Scale**: 4, 8, 12, 16, 20, 24, 32, 40, 48, 64, 80, 96px

## 🔌 API Integration

The app expects a REST API with the following endpoint:

```
GET /jobs?q={query}&limit={limit}&offset={offset}&experience_level={level}&employment_type={type}&location={location}&work_mode={mode}&company={company}
```

Response format:
```json
{
  "data": [
    {
      "job_id": 123,
      "company_id": 45,
      "company_name": "TechCorp",
      "company_logo_url": "https://example.com/logo.png",
      "title": "Software Engineer",
      "description": "Job description...",
      "experience_level": "Senior",
      "employment_type": "Full-time",
      "location": "Costa Rica",
      "work_mode": "Remote",
      "application_url": "https://example.com/apply",
      "technologies": [
        {
          "name": "React",
          "category": "Framework",
          "required": true
        }
      ],
      "posted_at": "2024-12-15T10:30:00Z"
    }
  ],
  "pagination": {
    "total": 150,
    "limit": 20,
    "offset": 0,
    "has_more": true
  }
}
```

## 🧪 Testing

Run tests with:
```bash
npm test
```

## 🚀 Building for Production

Build the app:
```bash
npm run build
```

The build files will be in the `build/` directory.

## 📈 Performance

- **Target load time**: < 2 seconds
- **Core Web Vitals**: LCP < 2.5s, FID < 100ms, CLS < 0.1
- **Image optimization**: WebP with fallbacks
- **Code splitting**: Ready for route-based splitting

## 🔮 Future Enhancements

### Phase 2
- [ ] User authentication
- [ ] Saved job management
- [ ] Email job alerts
- [ ] Company profiles

### Phase 3
- [ ] Job recommendations
- [ ] Advanced filters
- [ ] Salary insights
- [ ] Application tracking

### Phase 4
- [ ] Analytics dashboard
- [ ] A/B testing
- [ ] Performance optimizations
- [ ] Mobile app

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License.

## 📞 Support

For questions or support, contact: hello@ticosintech.com

---

Made with ❤️ in 🇨🇷 Costa Rica