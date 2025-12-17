package handler

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"web_service_auth/internal/app/ds"
	"web_service_auth/internal/app/role"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

const jwtPrefix = "Bearer "

func (h *Handler) WithAuthCheck(assignedRoles ...role.Role) func(ctx *gin.Context) {
	return func(gCtx *gin.Context) {
		jwtStr := gCtx.GetHeader("Authorization")

		if !strings.HasPrefix(jwtStr, jwtPrefix) { // если нет префикса то нас дурят!
			gCtx.AbortWithStatus(http.StatusForbidden) // отдаем что нет доступа
			return                                     // завершаем обработку
		}

		// отрезаем префикс
		jwtStr = jwtStr[len(jwtPrefix):]
		// проверяем jwt в блеклист редиса
		err := h.Redis.CheckJWTInBlacklist(gCtx.Request.Context(), jwtStr)
		if err == nil { // значит что токен в блеклисте
			gCtx.AbortWithStatus(http.StatusForbidden)

			return
		}
		if !errors.Is(err, redis.Nil) { // значит что это не ошибка отсуствия - внутренняя ошибка
			gCtx.AbortWithError(http.StatusInternalServerError, err)

			return
		}

		token, err := jwt.ParseWithClaims(jwtStr, &ds.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(h.Config.JWT.Token), nil
		})
		if err != nil {
			gCtx.AbortWithStatus(http.StatusForbidden)
			log.Println(err)

			return
		}

		myClaims := token.Claims.(*ds.JWTClaims)

		for _, oneOfAssignedRole := range assignedRoles {
			if myClaims.Role == oneOfAssignedRole {
				gCtx.Next()
				return
			}
		}
		gCtx.AbortWithStatus(http.StatusForbidden)
		log.Printf("role %s is not assigned in %s", myClaims.Role, assignedRoles)

		return

	}

}

func (h *Handler) GetCurrentUserID(ctx *gin.Context) uuid.UUID {
	jwtStr := ctx.GetHeader("Authorization")

	jwtStr = jwtStr[len(jwtPrefix):]

	token, _ := jwt.ParseWithClaims(jwtStr, &ds.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.Config.JWT.Token), nil
	})

	myClaims := token.Claims.(*ds.JWTClaims)
	return myClaims.UserUUID

}

func (h *Handler) GetCurrentUserRole(ctx *gin.Context) role.Role {
	jwtStr := ctx.GetHeader("Authorization")

	jwtStr = jwtStr[len(jwtPrefix):]

	token, _ := jwt.ParseWithClaims(jwtStr, &ds.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.Config.JWT.Token), nil
	})

	myClaims := token.Claims.(*ds.JWTClaims)
	return myClaims.Role

}
