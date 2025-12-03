import { FC, useEffect, useState } from 'react';
import { Container, Row, Col, Form, Button, Spinner } from 'react-bootstrap';
import { useNavigate } from 'react-router-dom';
import type { CarbonateInfo } from '../modules/types';
import { getCurrentCarbonate } from '../modules/api';
import { AcidCard } from '../components/AcidCard';
import { BreadCrumbs } from '../components/BreadCrumbs';
import { ROUTE_LABELS } from '../Routes';
import './AcidsPage.css';

// Redux imports
import { useDispatch } from 'react-redux';
import {
    useAcidsList,
    useSearchTerm,
    setSearchTermAction
} from '../slices/acidsSlice';
import { UseAcidsData } from '../hooks/useAcidsData';

export const AcidsPage: FC = () => {
    const navigate = useNavigate();
    const dispatch = useDispatch();

    // Получаем данные из Redux через селекторы
    const acids = useAcidsList();
    const reduxSearch = useSearchTerm();

    // Локальное состояние для input поля, чтобы не дергать Redux на каждый символ
    // (обновляем Redux только при Submit)
    const [localSearch, setLocalSearch] = useState(reduxSearch);

    // Карбонаты пока оставим локально, так как задача стояла про поиск и кислоты
    const [carbonateInfo, setCarbonateInfo] = useState<CarbonateInfo | null>(null);

    // Подключаем наш хук для загрузки данных.
    // Он сам следит за изменением reduxSearch и обновляет список acids.
    const { loading: acidsLoading } = UseAcidsData();

    // Загрузка карбонатов отдельно (так как это не часть задачи по Redux)
    useEffect(() => {
        getCurrentCarbonate().then(setCarbonateInfo);
    }, []);

    // Синхронизируем локальное поле ввода, если в Redux изменился поиск извне
    useEffect(() => {
        setLocalSearch(reduxSearch);
    }, [reduxSearch]);

    const handleSearchSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        // Отправляем действие в Redux. Это триггернет useEffect в UseAcidsData
        dispatch(setSearchTermAction(localSearch));
    };

    const isCalcActive = carbonateInfo && carbonateInfo.AcidCount > 0;

    return (
        <Container>
            <BreadCrumbs crumbs={[{ label: ROUTE_LABELS.ACIDS }]} />

            <div className="top-row-container">
                <Button className="calc-button-style" disabled={!isCalcActive} onClick={() => isCalcActive && navigate(`/carbonate/${carbonateInfo?.CarbonateID}`)}>
                    🧪 Рассчитать{isCalcActive && (<span className="indicator">{carbonateInfo?.AcidCount}</span>)}
                </Button>

                <Form className="flex-grow-1" onSubmit={handleSearchSubmit}>
                    <Form.Control
                        type="search"
                        placeholder="Найти кислоту..."
                        value={localSearch}
                        onChange={(e) => setLocalSearch(e.target.value)}
                    />
                </Form>
            </div>

            {acidsLoading ? (
                <div className="d-flex justify-content-center p-5">
                    <Spinner animation="border" style={{ color: '#004976' }} />
                </div>
            ) : (
                <>
                    {acids.length === 0 ? (
                        <p className="text-center text-muted">Ничего не найдено.</p>
                    ) : (
                        <Row xs={1} md={2} lg={3} className="g-4">
                            {acids.map((acid) => (
                                <Col key={acid.ID}>
                                    <AcidCard
                                        id={acid.ID}
                                        name={acid.Name}
                                        img={acid.Img}
                                    />
                                </Col>
                            ))}
                        </Row>
                    )}
                </>
            )}
        </Container>
    );
};