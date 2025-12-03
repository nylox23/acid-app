import type { FC } from "react";
import { Container, Button } from "react-bootstrap";
import { Link } from "react-router-dom";
import { ROUTES } from "../Routes";

export const HomePage: FC = () => {
    return (
        <Container className="p-5 mb-4 bg-light rounded-3">
            <h1>Калькулятор Карбонатов</h1>
            <p className="lead">
                Это приложение позволяет просматривать информацию о различных кислотах,
                узнавать их молярную массу и количество ионов водорода и проводить расчет
                количества газа, выделяющегося при взаимодействии CaCO3 с ними (WIP)
            </p>
            <hr className="my-4" />
            <p>
                Перейдите к списку, чтобы найти нужную кислоту.
            </p>
            <Link to={ROUTES.ACIDS}>
                <Button variant="primary" size="lg">Список кислот</Button>
            </Link>
        </Container>
    );
};