import { BrowserRouter, Route, Routes } from "react-router-dom";
import MyNavbar from "./components/MyNavbar.tsx";
import { HomePage } from "./pages/HomePage";
import { AcidsPage } from "./pages/AcidsPage";
import { AcidDetailsPage } from "./pages/AcidDetailsPage";
import { ROUTES } from "./Routes";
import Container from "react-bootstrap/Container";

// То же имя, что и в vite.config.ts
const REPO_NAME = '/acid-app';

function App() {
    return (
        // Добавляем basename
        <BrowserRouter basename={REPO_NAME}>
            <MyNavbar />
            <Container className="mt-4">
                <Routes>
                    <Route path={ROUTES.HOME} element={<HomePage />} />
                    <Route path={ROUTES.ACIDS} element={<AcidsPage />} />
                    <Route path={`${ROUTES.ACIDS}/:id`} element={<AcidDetailsPage />} />
                </Routes>
            </Container>
        </BrowserRouter>
    );
}

export default App;