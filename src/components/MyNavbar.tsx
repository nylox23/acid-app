import Container from 'react-bootstrap/Container';
import Nav from 'react-bootstrap/Nav';
import Navbar from 'react-bootstrap/Navbar';
import Button from 'react-bootstrap/Button';
import { Link, useNavigate } from 'react-router-dom';
import { ROUTES, ROUTE_LABELS } from '../Routes';
import { useSelector, useDispatch } from 'react-redux';
import type { RootState, AppDispatch } from '../store';
import { logoutUser } from '../slices/authSlice';
import { useEffect } from 'react';
import { getProfile } from '../slices/authSlice';
import { fetchCurrentCarbonateInfo } from '../slices/carbonateSlice';

function MyNavbar() {
    const dispatch = useDispatch<AppDispatch>();
    const navigate = useNavigate();
    const { isAuthenticated, user } = useSelector((state: RootState) => state.auth);
    const { currentCarbonateId, currentAcidCount } = useSelector((state: RootState) => state.carbonateData);

    useEffect(() => {
        if (isAuthenticated) {
            dispatch(getProfile());
            dispatch(fetchCurrentCarbonateInfo());
        }
    }, [isAuthenticated, dispatch]);

    const handleLogout = () => {
        dispatch(logoutUser());
        navigate(ROUTES.LOGIN);
    };

    return (
        <Navbar className="navbar-custom" variant="dark" expand="lg">
            <Container>
                <Navbar.Brand as={Link} to={ROUTES.HOME} className="btn-home me-3">
                    ⌂ Домой
                </Navbar.Brand>

                <Navbar.Toggle aria-controls="basic-navbar-nav" />
                <Navbar.Collapse id="basic-navbar-nav">
                    <Nav className="me-auto">
                        <Nav.Link as={Link} to={ROUTES.ACIDS}>{ROUTE_LABELS.ACIDS}</Nav.Link>
                        {isAuthenticated && (
                            <Nav.Link as={Link} to={ROUTES.CARBONATE_LIST}>{ROUTE_LABELS.CARBONATE_LIST}</Nav.Link>
                        )}
                    </Nav>
                    <Nav className="align-items-center gap-3">
                        {isAuthenticated && (
                            <Button
                                variant={currentCarbonateId ? "primary" : "secondary"}
                                disabled={!currentCarbonateId}
                                onClick={() => currentCarbonateId && navigate(`${ROUTES.CARBONATE_DETAIL}/${currentCarbonateId}`)}
                                className="position-relative"
                            >
                                Калькулятор
                                {currentAcidCount > 0 && (
                                    <span className="position-absolute top-0 start-100 translate-middle badge rounded-pill bg-danger">
                                        {currentAcidCount}
                                    </span>
                                )}
                            </Button>
                        )}

                        {isAuthenticated ? (
                            <>
                                <Navbar.Text className="text-white">
                                    Привет, <Link to={ROUTES.PROFILE} className="text-white fw-bold" style={{textDecoration: 'none'}}>{user?.login}</Link>
                                </Navbar.Text>
                                <Button variant="outline-light" size="sm" onClick={handleLogout}>Выход</Button>
                            </>
                        ) : (
                            <Nav.Link as={Link} to={ROUTES.LOGIN}>Вход</Nav.Link>
                        )}
                    </Nav>
                </Navbar.Collapse>
            </Container>
        </Navbar>
    );
}

export default MyNavbar;