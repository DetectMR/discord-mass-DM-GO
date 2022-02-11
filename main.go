// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package main

import (
	"bufio"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/V4NSH4J/discord-mass-dm-GO/utilities"
	"github.com/fatih/color"
	"github.com/zenthangplus/goccm"
)

func main() {
	// Credits
	CaptchaServices = []string{"capmonster.cloud", "anti-captcha.com", ""}
	color.Blue(logo)
	color.Green("Made by https://github.com/V4NSH4J\nStar repository on github for updates!")
	Options()
}

// Options menu
func Options() {
	reg := regexp.MustCompile(`(.+):(.+):(.+)`)
	color.Green("Меню:\n |- 01) Инвайт на сервер [для Токенов]\n |- 02) Массовая рассылка в ЛС [для Токенов]\n |- 03) Одиночная рассылка в ЛС [для Токенов]\n |- 04) Накрутка реакций  [для Токенов]\n |- 05) Получить сообщение[Input]\n |- 06) Конвертировать mail:Pass:Token в Token [Email:Password:Token]\n |- 07) Токен чекер [для Токенов]]\n |- 08) Выход из сервера [для Токенов]]\n |- 09) Поддержка токенов онлайн[для Токенов]\n |- 10) Парсинг Меню [Input]\n |- 11) Изменить имя [Email:Password:Token]\n |- 12) Изменить аватарку [для Токенов]\n |- 13) Проверьте, находятся ли ваши токены на сервере [для Токенов]\n |- 14) Информация\n |- 15) Выход")
	color.Red("\nВведите свой выбор: ")
	var choice int
	fmt.Scanln(&choice)
	switch choice {
	default:
		color.Red("Неверный выбор!")
		Options()
	case 0:
		color.Cyan("Режим отладки")
	case 1:
		var invitechoice int
		color.White("color.Green("Инвайт Меню:\n1) Одиночный инвайт\n2) Масовый инвайт из файла")
")
		fmt.Scanln(&invitechoice)
		if invitechoice != 1 && invitechoice != 2 {
			color.Red("[%v] Invalid choice", time.Now().Format("15:04:05"))
			ExitSafely()
			return
		}
		switch invitechoice {
		case 1:
			color.Cyan("режим одиночного инвайта")
			color.Green("Это позволит присоединить ваши токены из файла tokens.txt к серверу")
			cfg, instances, err := getEverything()
			if err != nil {
				color.Red("[%v] Ошибка при получении необходимых данных: %v", time.Now().Format("15:04:05"), err)
			}
			color.Green("[%v] Введите приглашения на сервер (всё после discord.gg/): ", time.Now().Format("15:04:05"))
			var invite string
			fmt.Scanln(&invite)
			color.Green("[%v] Введите количество потоков (0: Неограниченное количество. 1: Для использования надлежащей задержки): ", time.Now().Format("15:04:05"))
			var threads int
			fmt.Scanln(&threads)

			if threads > len(instances) {
				threads = len(instances)
			}
			if threads == 0 {
				threads = len(instances)
			}

			color.Green("[%v] Введите базовую задержку для присоединения в секундах (0 - нет)", time.Now().Format("15:04:05"))
			var base int
			fmt.Scanln(&base)
			color.Green("[%v] Введите случайную задержку, которая будет добавлена к базовой задержке (0 - нет).", time.Now().Format("15:04:05"))
			var random int
			fmt.Scanln(&random)
			var delay int
			if random > 0 {
				delay = base + rand.Intn(random)
			} else {
				delay = base
			}
			c := goccm.New(threads)
			for i := 0; i < len(instances); i++ {
				time.Sleep(time.Duration(cfg.Offset) * time.Millisecond)
				c.Wait()
				go func(i int) {
					err := instances[i].Invite(invite)
					if err != nil {
						color.Red("[%v] Ошибка при присоединении: %v", time.Now().Format("15:04:05"), err)
					}
					time.Sleep(time.Duration(delay) * time.Second)
					c.Done()

				}(i)
			}
			c.WaitAllDone()
			color.Green("[%v] Все потоки закончены", time.Now().Format("15:04:05"))

		case 2:
			color.Cyan("Режим массового инвайта")
			color.Green("Это соединит ваши токены из файла tokens.txt с серверами из файла invite.txt")
			cfg, instances, err := getEverything()
			if err != nil {
				color.Red("[%v] Ошибка при получении необходимых данных: %v", time.Now().Format("15:04:05"), err)
			}

			if len(instances) == 0 {
				color.Red("[%v] Введите свои токены в файл tokens.txt", time.Now().Format("15:04:05"))
				ExitSafely()
			}
			invites, err := utilities.ReadLines("invite.txt")
			if err != nil {
				color.Red("Ошибка при открытии файла invite.txt: %v", err)
				ExitSafely()
				return
			}
			if len(invites) == 0 {
				color.Red("[%v] Введите свои приглашения в файл invite.txt", time.Now().Format("15:04:05"))
				ExitSafely()
				return
			}
			color.Green("Введите задержку между 2 последовательными присоединениями по 1 Токену в секундах: ")
			var delay int
			fmt.Scanln(&delay)
			color.Green("Введите количество потоков (0 для неограниченного количества): ")
			var threads int
			fmt.Scanln(&threads)
			if threads > len(instances) {
				threads = len(instances)
			}
			if threads == 0 {
				threads = len(instances)
			}
			c := goccm.New(threads)
			for i := 0; i < len(instances); i++ {
				time.Sleep(time.Duration(cfg.Offset) * time.Millisecond)
				c.Wait()
				go func(i int) {
					for j := 0; j < len(invites); j++ {
						err := instances[i].Invite(invites[j])
						if err != nil {
							color.Red("[%v] Ошибка при присоединении: %v", time.Now().Format("15:04:05"), err)
						}
						time.Sleep(time.Duration(delay) * time.Second)
					}
					c.Done()
				}(i)
			}
			c.WaitAllDone()
			color.Green("[%v] Все потоки закончены", time.Now().Format("15:04:05"))
		}
	case 2:

		color.Cyan("Массовая рассылка В ЛС")
		color.Green("Отправка Личного Сообщения всем пользователям в файле memberids.txt  всеми вашими токенами в файле tokens.txt")
		members, err := utilities.ReadLines("memberids.txt")
		if err != nil {
			color.Red("Ошибка при открытии файла memberids.txt: %v", err)
			ExitSafely()
		}
		cfg, instances, err := getEverything()
		if err != nil {
			color.Red("[%v] Ошибка при получении необходимых данных: %v", time.Now().Format("15:04:05"), err)
		}
		var msg utilities.Message
		color.Green("Нажмите - (1) чтобы использовать сообщения из файла (message.json).\nНажмите - (2) чтобы ввести сообщение здесь: ")
		var messagechoice int
		fmt.Scanln(&messagechoice)
		if messagechoice != 1 && messagechoice != 2 {
			color.Red("[%v] Неверный выбор", time.Now().Format("15:04:05"))
			ExitSafely()
		}
		if messagechoice == 2 {
			color.Green("Введите свое сообщение, для новой строки исполйьзуйте код (\n). Вы также можете изпользовать этот код (\n)  в файле message.json")
			scanner := bufio.NewScanner(os.Stdin)
			var text string
			if scanner.Scan() {
				text = scanner.Text()
			}

			msg.Content = text
			msg.Content = strings.Replace(msg.Content, "\\n", "\n", -1)
			var msgs []utilities.Message
			msgs = append(msgs, msg)
			err := setMessages(instances, msgs)
			if err != nil {
				color.Red("[%v] Ошибка при вводе сообщения: %v", time.Now().Format("15:04:05"), err)
				ExitSafely()
			}
		} else {
			var msgs []utilities.Message
			err := setMessages(instances, msgs)
			if err != nil {
				color.Red("[%v] Ошибка при вводе сообщения: %v", time.Now().Format("15:04:05"), err)
				ExitSafely()
			}
		}
		color.Green("[%v] Хотите ли вы использовать расширенные настройки? 0: Нет, 1: Да:", time.Now().Format("15:04:05"))
		var advancedchoice int
		var checkchoice int
		var serverid string
		var tryjoinchoice int
		var invite string
		var maxattempts int
		fmt.Scanln(&advancedchoice)
		if advancedchoice != 0 && advancedchoice != 1 {
			color.Red("[%v] Неверный выбор", time.Now().Format("15:04:05"))
			ExitSafely()
		}
		if advancedchoice == 1 {
			color.White("[%v] Вы хотите проверять, находится ли токен на сервере перед каждым DM? [0: Нет, 1: Да]", time.Now().Format("15:04:05"))
			fmt.Scanln(&checkchoice)
			if checkchoice != 0 && checkchoice != 1 {
				color.Red("[%v] Неверный выбор", time.Now().Format("15:04:05"))
				ExitSafely()
			}
			if checkchoice == 1 {
				color.White("[%v] Введите ID сервера", time.Now().Format("15:04:05"))
				fmt.Scanln(&serverid)
				color.White("[%v] Хотите ли вы попробовать снова подключиться к серверу, если токен не находится на сервере? [0: Нет, 1: Да]", time.Now().Format("15:04:05"))
				fmt.Scanln(&tryjoinchoice)
				if tryjoinchoice != 0 && tryjoinchoice != 1 {
					color.Red("[%v] Неверный выбор", time.Now().Format("15:04:05"))
					ExitSafely()
				}
				if tryjoinchoice == 1 {
					color.White("[%v] Введите постоянный код приглашения", time.Now().Format("15:04:05"))
					fmt.Scanln(&invite)
					color.White("[%v] Введите максимальное количество попыток повторного присоединения", time.Now().Format("15:04:05"))
					fmt.Scanln(&maxattempts)
				}
			}
		}
		// Also initiate variables and slices for logging and counting
		var session []string
		var completed []string
		var failed []string
		var dead []string
		var failedCount = 0
		completed, err = utilities.ReadLines("completed.txt")
		if err != nil {
			color.Red("Ошибка при открытии файла completed.txt: %v", err)
			ExitSafely()
		}
		if cfg.Skip {
			members = utilities.RemoveSubset(members, completed)
		}
		if cfg.SkipFailed {
			failedSkip, err := utilities.ReadLines("failed.txt")
			if err != nil {
				color.Red("Ошибка при открытии файла failed.txt: %v", err)
				ExitSafely()
			}
			members = utilities.RemoveSubset(members, failedSkip)
		}
		if len(instances) == 0 {
			color.Red("[%v] Введите свои токены в файл tokens.txt ", time.Now().Format("15:04:05"))
			ExitSafely()
		}
		if len(members) == 0 {
			color.Red("[%v] Введите идентификаторы участников в файл memberids.txt или убедитесь, что все они не находятся в файле completed.txt", time.Now().Format("15:04:05"))
			ExitSafely()
		}
		if len(members) < len(instances) {
			instances = instances[:len(members)]
		}
		msgs := instances[0].Messages
		for i := 0; i < len(msgs); i++ {
			if msgs[i].Content == "" && msgs[i].Embeds == nil {
				color.Red("[%v] ПРЕДУПРЕЖДЕНИЕ: Сообщение %v пустое", time.Now().Format("15:04:05"), i)
			}
		}
		// Send members to a channel
		mem := make(chan string, len(members))
		go func() {
			for i := 0; i < len(members); i++ {
				mem <- members[i]
			}
		}()
		// Setting information to windows titlebar by github.com/foxzsz
		go func() {
			for {
				cmd := exec.Command("cmd", "/C", "title", fmt.Sprintf(`DMDGO [%d sent, %v failed, %d locked, %d avg. dms, %d tokens left]`, len(session), len(failed), len(dead), len(session)/len(instances), len(instances)-len(dead)))
				_ = cmd.Run()
			}
		}()
		var wg sync.WaitGroup
		start := time.Now()
		for i := 0; i < len(instances); i++ {
			// Offset goroutines by a few milliseconds. Makes a big difference and allows for better concurrency
			time.Sleep(time.Duration(cfg.Offset) * time.Millisecond)
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				for {
					// Get a member from the channel
					if len(mem) == 0 {
						break
					}
					member := <-mem

					// Breaking loop if maximum DMs reached
					if cfg.MaxDMS != 0 && instances[i].Count >= cfg.MaxDMS {
							color.Yellow("[%v] Максимальные значения ЛС достигнуты для %v", time.Now().Format("15:04:05"), instances[i].Token)
						break
					}
					// Start websocket connection if not already connected and reconnect if dead
					if cfg.Websocket && instances[i].Ws == nil {
						err := instances[i].StartWS()
						if err != nil {
							color.Red("[%v] Ошибка при открытии вебсокета: %v", time.Now().Format("15:04:05"), err)
						} else {
							color.Green("[%v] Открыт вебсокет %v", time.Now().Format("15:04:05"), instances[i].Token)
						}
					}
					if cfg.Websocket && cfg.Receive && instances[i].Ws != nil && !instances[i].Receiver {
						instances[i].Receiver = true
						go func() {
							for {
								if !instances[i].Receiver {
									break
								}
								mes := <-instances[i].Ws.Messages
								if !strings.Contains(string(mes), "guild_id") {
									var mar utilities.Event
									err := json.Unmarshal(mes, &mar)
									if err != nil {
										color.Red("[%v] Error while unmarshalling websocket message: %v", time.Now().Format("15:04:05"), err)
										continue
									}
									if instances[i].ID == "" {
										tokenPart := strings.Split(instances[i].Token, ".")[0]
										dec, err := base64.StdEncoding.DecodeString(tokenPart)
										if err != nil {
											color.Red("[%v] Ошибка при декодировании токена: %v", time.Now().Format("15:04:05"), err)
											continue
										}
										instances[i].ID = string(dec)
									}
									if mar.Data.Author.ID == instances[i].ID {
										continue
									}
									color.Green("[%v] %v#%v отправил сообщение %v : %v", time.Now().Format("15:04:05"), mar.Data.Author.Username, mar.Data.Author.Discriminator, instances[i].Token, mar.Data.Content)
									newStr := "Username: " + mar.Data.Author.Username + "#" + mar.Data.Author.Discriminator + "\nID: " + mar.Data.Author.ID + "\n" + "Message: " + mar.Data.Content + "\n"
									err = utilities.WriteLines("received.txt", newStr)
									if err != nil {
										color.Red("[%v] Ошибка при открытии файла received.txt: %v", time.Now().Format("15:04:05"), err)
									}
								}
							}
						}()
					}
					// Check if token is valid
					status := instances[i].CheckToken()
					if status != 200 && status != 204 && status != 429 && status != -1 {
						failedCount++
						color.Red("[%v] Токен %v может быть заблокирован - Остановка потока и добавление пользователя в файл failed.txt. %v [%v]", time.Now().Format("15:04:05"), instances[i].Token, status, failedCount)
						failed = append(failed, member)
						dead = append(dead, instances[i].Token)
						err := utilities.WriteLines("failed.txt", member)
						if err != nil {
							fmt.Println(err)
						}
						if cfg.Stop {
							break
						}
					}
					// Advanced Options
					if advancedchoice == 1 {
						if checkchoice == 1 {
							r, err := instances[i].ServerCheck(serverid)
							if err != nil {
								color.Red("[%v] Ошибка при проверке сервера: %v", time.Now().Format("15:04:05"), err)
								continue
							}
							if r != 200 && r != 204 && r != 429 {
								if tryjoinchoice == 0 {
									color.Red("[%v] Stopping token %v [Нет на сервере]", time.Now().Format("15:04:05"), instances[i].Token)

									break
								} else {
									if instances[i].Rejoin >= maxattempts {
										color.Red("[%v] Stopping token %v [Максимальное количество попыток повторного подключения к серверу]", time.Now().Format("15:04:05"), instances[i].Token)
										break
									}
									err := instances[i].Invite(invite)
									if err != nil {
										color.Red("[%v] Ошибка при подключении к серверу: %v", time.Now().Format("15:04:05"), err)
										instances[i].Rejoin++
										continue
									}
								}
							}
						}
					}
					var user string
					user = member
					// Check Mutual
					if cfg.Mutual {
						info, err := instances[i].UserInfo(member)
						if err != nil {
							failedCount++
							color.Red("[%v] Ошибка при получении информации о пользователе: %v [%v]", time.Now().Format("15:04:05"), err, failedCount)
							err = WriteLine("input/failed.txt", member)
							if err != nil {
								fmt.Println(err)
							}
							failed = append(failed, member)

							continue
						}
						if len(info.Mutual) == 0 {
							failedCount++
							color.Red("[%v] Токен %v не смог отправить сообщения в ЛС %v [Нет общих серверов] [%v]", time.Now().Format("15:04:05"), instances[i].Token, info.User.Username+info.User.Discriminator, failedCount)
							err = WriteLine("input/failed.txt", member)
							if err != nil {
								fmt.Println(err)
							}
							failed = append(failed, member)
							continue
						}
						user = info.User.Username + "#" + info.User.Discriminator
						// Used only if Websocket is enabled as Unwebsocketed Tokens get locked if they attempt to send friend requests.
						if cfg.Friend && cfg.Websocket {
							x, err := strconv.Atoi(info.User.Discriminator)
							if err != nil {
								color.Red("[%v] Ошибка при добавлении друга: %v", time.Now().Format("15:04:05"), err)
							}
							resp, err := instances[i].Friend(info.User.Username, x)
							if err != nil {
								color.Red("[%v] Ошибка при добавлении друга: %v", time.Now().Format("15:04:05"), err)
							}
							if resp.StatusCode != 204 && err != nil {
								body, _ := utilities.ReadBody(*resp)
								color.Red("[%v] Ошибка при добавлении друга: %v", time.Now().Format("15:04:05"), string(body))
							} else {
								color.Green("[%v] Добавленный друг %v", time.Now().Format("15:04:05"), info.User.Username+"#"+info.User.Discriminator)
							}
						}
					}
					// Open channel to get snowflake
					snowflake, err := instances[i].OpenChannel(member)
					if err != nil {
						failedCount++
						color.Red("[%v] Ошибка при открытии  ЛС: %v [%v]", time.Now().Format("15:04:05"), err, failedCount)
						err = WriteLine("input/failed.txt", member)
						if err != nil {
							fmt.Println(err)
						}
						failed = append(failed, member)
						continue
					}
					resp, err := instances[i].SendMessage(snowflake, member)
					if err != nil {
						failedCount++
						color.Red("[%v] Ошибка при отправке сообщения: %v [%v]", time.Now().Format("15:04:05"), err, failedCount)
						err = WriteLine("input/failed.txt", member)
						if err != nil {
							fmt.Println(err)
						}
						failed = append(failed, member)
						continue
					}
					body, err := utilities.ReadBody(resp)
					if err != nil {
						failedCount++
						color.Red("[%v] Ошибка при чтении тела: %v [%v]", time.Now().Format("15:04:05"), err, failedCount)
						err = WriteLine("input/failed.txt", member)
						if err != nil {
							fmt.Println(err)
						}
						failed = append(failed, member)
						continue
					}
					var response jsonResponse
					errx := json.Unmarshal(body, &response)
					if errx != nil {
						failedCount++
						color.Red("[%v] Error while unmarshalling body: %v [%v]", time.Now().Format("15:04:05"), errx, failedCount)
						err = WriteLine("input/failed.txt", member)
						if err != nil {
							fmt.Println(err)
						}
						failed = append(failed, member)
						continue
					}
					// Everything is fine, continue as usual
					if resp.StatusCode == 200 {
						err = WriteLine("input/completed.txt", member)
						if err != nil {
							fmt.Println(err)
						}
						completed = append(completed, member)
						session = append(session, member)
						color.Green("[%v] Токен %v отправил ЛС на %v [%v]", time.Now().Format("15:04:05"), instances[i].Token, user, len(session))
						if cfg.Websocket && cfg.Call && instances[i].Ws != nil {
							err := instances[i].Call(snowflake)
							if err != nil {
								color.Red("[%v] %v Ошибка при вызове %v: %v", time.Now().Format("15:04:05"), instances[i].Token, user, err)
							}
							// Unfriended people can't ring.
							//
							// resp, err := utilities.Ring(instances[i].Client, instances[i].Token, snowflake)
							// if err != nil {
							//      color.Red("[%v] %v Error while ringing %v: %v", time.Now().Format("15:04:05"), instances[i].Token, user, err)
							// }
							// if resp == 200 || resp == 204 {
							//      color.Green("[%v] %v Ringed %v", time.Now().Format("15:04:05"), instances[i].Token, user)
							// } else {
							//      color.Red("[%v] %v Error while ringing %v: %v", time.Now().Format("15:04:05"), instances[i].Token, user, resp)
							// }

						}
						// Forbidden - Token is being rate limited
					} else if resp.StatusCode == 403 && response.Code == 40003 {

						err = WriteLine("input/failed.txt", member)
						if err != nil {
							fmt.Println(err)
						}
						color.Yellow("[%v] Токен %v в тайм ауте %v минут!", time.Now().Format("15:04:05"), instances[i].Token, int(cfg.LongDelay/60))
						time.Sleep(time.Duration(cfg.LongDelay) * time.Second)
						color.Yellow("[%v] Токен %v продолжил!", time.Now().Format("15:04:05"), instances[i].Token)
						// Forbidden - DM's are closed
					} else if resp.StatusCode == 403 && response.Code == 50007 {
						failedCount++
						err = WriteLine("input/failed.txt", member)
						if err != nil {
							fmt.Println(err)
						}
						color.Red("[%v] Токен %v не удалось отправить ЛС %v Пользователь имеет закрытое ЛС или отсутствует на сервере %v [%v]", time.Now().Format("15:04:05"), instances[i].Token, user, string(body), failedCount)
						// Forbidden - Locked or Disabled
					} else if (resp.StatusCode == 403 && response.Code == 40002) || resp.StatusCode == 401 || resp.StatusCode == 405 {
						failedCount++
						err = WriteLine("input/failed.txt", member)
						if err != nil {
							fmt.Println(err)
						}
						color.Red("[%v] Токен %v заблокирован или отключен. Остановка экземпляра. %v %v [%v]", time.Now().Format("15:04:05"), instances[i].Token, resp.StatusCode, string(body), failedCount)
						dead = append(dead, instances[i].Token)
						// Stop token if locked or disabled
						if cfg.Stop {
							break
						}
						// Forbidden - Invalid token
					} else if resp.StatusCode == 403 && response.Code == 50009 {
						failedCount++
						err = WriteLine("input/failed.txt", member)
						if err != nil {
							fmt.Println(err)
						}
						color.Red("[%v] Токен %v не может отправить ЛС %v. Возможно, он не прошел проверку на членство, или его уровень проверки слишком низок, или сервер требует от новых членов подождать 10 минут, прежде чем они смогут взаимодействовать на сервере.%v [%v]", time.Now().Format("15:04:05"), instances[i].Token, user, string(body), failedCount)
						// General case - Continue loop. If problem with instance, it will be stopped at start of loop.
					} else if resp.StatusCode == 429 {
						color.Red("[%v] Токен %v ограничен лимитом. Сон в течение 10 секунд", time.Now().Format("15:04:05"), instances[i].Token)
						time.Sleep(10 * time.Second)
					} else {
						failedCount++
						err = WriteLine("input/failed.txt", member)
						if err != nil {
							fmt.Println(err)
						}
						color.Red("[%v] Токен %v не смог отправить ЛС %v Код ошибки: %v; Статус: %v; Сообщение: %v [%v]", time.Now().Format("15:04:05"), instances[i].Token, user, response.Code, resp.Status, response.Message, failedCount)
					}
					time.Sleep(time.Duration(cfg.Delay) * time.Second)
				}
			}(i)
		}
		wg.Wait()

		color.Green("[%v] Потоки завершени! Запись в файл", time.Now().Format("15:04:05"))

		elapsed := time.Since(start)
		color.Green("[%v] Рассылка ЛС занимала %v. Успешно отправлены ЛС на %v ID. Не удалось отправить ЛС на %v ID. %v Токены не функционируют и %v Токены функционируют", time.Now().Format("15:04:05"), elapsed.Seconds(), len(completed), len(failed), len(dead), len(instances)-len(dead))
		if cfg.Remove {
			var tokens []string
			for i := 0; i < len(instances); i++ {
				tokens = append(tokens, instances[i].Token)
			}
			m := utilities.RemoveSubset(tokens, dead)
			err := Truncate("input/tokens.txt", m)
			if err != nil {
				fmt.Println(err)
			}
			color.Green("Обновленния файла tokens.txt")
		}
		if cfg.RemoveM {
			m := utilities.RemoveSubset(members, completed)
			err := Truncate("input/memberids.txt", m)
			if err != nil {
				fmt.Println(err)
			}
			color.Green("Updated memberids.txt")

		}
		if cfg.Websocket {
			for i := 0; i < len(instances); i++ {
				if instances[i].Ws != nil {
					instances[i].Ws.Close()
				}
			}
		}

	case 3:
		color.Cyan("Одиночный спамер")
		color.White("Введите 0 для одного сообщения; Введите 1 для непрерывного спама")
		var choice int
		fmt.Scanln(&choice)
		cfg, instances, err := getEverything()
		if err != nil {
			fmt.Println(err)
			ExitSafely()
		}
		var msg utilities.Message
		color.White("Нажмите 1, чтобы использовать сообщение из файла, или нажмите 2, чтобы ввести сообщение: ")
		var messagechoice int
		fmt.Scanln(&messagechoice)
		if messagechoice != 1 && messagechoice != 2 {
			color.Red("[%v] Неверный выбор", time.Now().Format("15:04:05"))
			ExitSafely()
		}
		if messagechoice == 2 {
			color.White("Введите свое сообщение, для новой строки используйте \\n. Чтобы использовать вставку, поместите сообщение в файл message.json: ")
			scanner := bufio.NewScanner(os.Stdin)
			var text string
			if scanner.Scan() {
				text = scanner.Text()
			}

			msg.Content = text
			msg.Content = strings.Replace(msg.Content, "\\n", "\n", -1)
			var msgs []utilities.Message
			msgs = append(msgs, msg)
			err := setMessages(instances, msgs)
			if err != nil {
				color.Red("[%v] Ошибка в сообщений: %v", time.Now().Format("15:04:05"), err)
				ExitSafely()
			}
		} else {
			var msgs []utilities.Message
			err := setMessages(instances, msgs)
			if err != nil {
				color.Red("[%v] Ошибка в сообщений: %v", time.Now().Format("15:04:05"), err)
				ExitSafely()
			}
		}

		color.White("Обеспечьте общую связь и введите ID жертвы: ")
		var victim string
		fmt.Scanln(&victim)
		var wg sync.WaitGroup
		wg.Add(len(instances))
		if choice == 0 {
			for i := 0; i < len(instances); i++ {
				time.Sleep(time.Duration(cfg.Offset) * time.Millisecond)

				go func(i int) {
					defer wg.Done()
					snowflake, err := instances[i].OpenChannel(victim)
					if err != nil {
						fmt.Println(err)
					}
					resp, err := instances[i].SendMessage(snowflake, victim)
					if err != nil {
						fmt.Println(err)
					}
					body, err := utilities.ReadBody(resp)
					if err != nil {
						fmt.Println(err)
					}
					if resp.StatusCode == 200 {
						color.Green("[%v] Токен %v ЛС %v", time.Now().Format("15:04:05"), instances[i].Token, victim)
					} else {
						color.Red("[%v] Токен %v ошибка при отправки ЛС %v [%v]", time.Now().Format("15:04:05"), instances[i].Token, victim, string(body))
					}
				}(i)
			}
			wg.Wait()
		}
		if choice == 1 {
			for i := 0; i < len(instances); i++ {
				time.Sleep(time.Duration(cfg.Offset) * time.Millisecond)
				go func(i int) {
					defer wg.Done()

					var c int
					for {
						snowflake, err := instances[i].OpenChannel(victim)
						if err != nil {
							fmt.Println(err)
						}
						resp, err := instances[i].SendMessage(snowflake, victim)
						if err != nil {
							fmt.Println(err)
						}
						if resp.StatusCode == 200 {
							color.Green("[%v] Токен %v ЛС %v [%v]", time.Now().Format("15:04:05"), instances[i].Token, victim, c)
						} else {
							color.Red("[%v] Токен %v ошибка при отправки ЛС %v", time.Now().Format("15:04:05"), instances[i].Token, victim)
						}
						c++
					}
				}(i)
				wg.Wait()
			}
		}
		color.Green("[%v] Потоки завершены!", time.Now().Format("15:04:05"))

	case 4:
		color.Cyan("Накрутка реакций")
		color.White("Примечание: Вам не нужно делать это, чтобы отправлять ЛС на серверах.")
		color.White("Меню:\n1) Из сообщения\n2) Вручную")
		var choice int
		fmt.Scanln(&choice)
		cfg, instances, err := getEverything()
		if err != nil {
			fmt.Println(err)
			ExitSafely()
		}
		var wg sync.WaitGroup
		wg.Add(len(instances))
		if choice == 1 {
			color.Cyan("Введите Токен:")
			var token string
			fmt.Scanln(&token)
			color.White("Введите ID сообщения: ")
			var id string
			fmt.Scanln(&id)
			color.White("Введите ID канала: ")
			var channel string
			fmt.Scanln(&channel)
			msg, err := utilities.GetRxn(channel, id, token)
			if err != nil {
				fmt.Println(err)
			}
			color.White("Выберите Emoji")
			for i := 0; i < len(msg.Reactions); i++ {
				color.White("%v) %v %v", i, msg.Reactions[i].Emojis.Name, msg.Reactions[i].Count)
			}
			var emoji int
			var send string
			fmt.Scanln(&emoji)
			for i := 0; i < len(instances); i++ {
				time.Sleep(time.Duration(cfg.Offset) * time.Millisecond)
				go func(i int) {
					defer wg.Done()
					if msg.Reactions[emoji].Emojis.ID == "" {
						send = msg.Reactions[emoji].Emojis.Name

					} else if msg.Reactions[emoji].Emojis.ID != "" {
						send = msg.Reactions[emoji].Emojis.Name + ":" + msg.Reactions[emoji].Emojis.ID
					}
					err := instances[i].React(channel, id, send)
					if err != nil {
						fmt.Println(err)
						color.Red("[%v] %v не отреагировал", time.Now().Format("15:04:05"), instances[i].Token)
					} else {
						color.Green("[%v] %v отреагировал на эмодзи", time.Now().Format("15:04:05"), instances[i].Token)
					}

				}(i)
			}
			wg.Wait()
			color.Green("[%v] Завершил все потоки.", time.Now().Format("15:04:05"))
		}
		if choice == 2 {
			color.Cyan("Введите ID канала")
			var channel string
			fmt.Scanln(&channel)
			color.White("Введите ID сообщения")
			var id string
			fmt.Scanln(&id)
			color.Red("Если у вас есть сообщение, пожалуйста, используйте вариант 1. Если вы хотите добавить пользовательский эмодзи. Следуйте этим инструкциям, если вы этого не сделаете, ничего не получится.\n Если это эмодзи по умолчанию, который появляется на клавиатуре эмодзи, просто скопируйте его как ТЕКСТ, а не как он появляется в Discord с двоеточиями. Если это пользовательский эмодзи (Nitro emoji), введите его следующим образом -> name:emojiID Чтобы получить идентификатор эмодзи, скопируйте ссылку на эмодзи и скопируйте идентификатор эмодзи из URL.\n Если вы не будете следовать этому, он не будет работать. Не пытайтесь делать невозможные вещи, например, пытаться запустить нитрореакцию с помощью ненитросчета.")
			color.White("Enter emoji")
			var emoji string
			fmt.Scanln(&emoji)
			for i := 0; i < len(instances); i++ {
				time.Sleep(time.Duration(cfg.Offset) * time.Millisecond)
				go func(i int) {
					defer wg.Done()
					err := instances[i].React(channel, id, emoji)
					if err != nil {
						fmt.Println(err)
						color.Red("[%v] %v не отреагировал", time.Now().Format("15:04:05"), instances[i].Token)
					}
					color.Green("[%v] %v отреагировал на эмодзи", time.Now().Format("15:04:05"), instances[i].Token)
				}(i)
			}
			wg.Wait()
			color.Green("[%v] Завершил все потоки.", time.Now().Format("15:04:05"))
		}

	case 5:
		// Uses ?around & ?limit parameters to discord's REST API to get messages to get the exact message needed
		color.Cyan("Получить сообщение - Это позволит получить ответное сообщение из Discord, которое вы хотите отправить..")
		color.White("Enter your token: \n")
		var token string
		fmt.Scanln(&token)
		color.White("Введите ID канала: \n")
		var channelID string
		fmt.Scanln(&channelID)
		color.White("Введите ID сообщения: \n")
		var messageID string
		fmt.Scanln(&messageID)
		message, err := utilities.FindMessage(channelID, messageID, token)
		if err != nil {
			color.Red("Ошибка при поиске сообщения: %v", err)
			ExitSafely()
			return
		}
		color.Green("[%v] Сообщение: %v", time.Now().Format("15:04:05"), message)

	case 6:
		// Quick way to interconvert tokens from a popular format to the one this program supports.
		color.Cyan("Email:Password:Token в Token")
		Tokens, err := utilities.ReadLines("tokens.txt")
		if err != nil {
			color.Red("Ошибка открытия tokens.txt: %v", err)
			ExitSafely()
			return
		}
		if len(Tokens) == 0 {
			color.Red("[%v] Введите ваши Токены в tokens.txt", time.Now().Format("15:04:05"))
			ExitSafely()
			return
		}
		var onlytokens []string
		for i := 0; i < len(Tokens); i++ {
			if strings.Contains(Tokens[i], ":") {
				token := strings.Split(Tokens[i], ":")[2]
				onlytokens = append(onlytokens, token)
			}
		}
		t := utilities.TruncateLines("tokens.txt", onlytokens)
		if t != nil {
			color.Red("[%v]При конвертирование ошибки tokens.txt: %v", time.Now().Format("15:04:05"), t)
			ExitSafely()
			return
		}

	case 7:
		// Basic token checker
		color.Cyan("Токен Чекер")
		cfg, instances, err := getEverything()
		if err != nil {
			color.Red("[%v] Ошибка при получении необходимых данных: %v", time.Now().Format("15:04:05"), err)
			ExitSafely()
		}
		color.White("Введите количество потоков: (0 для неограниченного количества)\n")
		var threads int
		fmt.Scanln(&threads)
		if threads > len(instances) {
			threads = len(instances)
		}
		if threads == 0 {
			threads = len(instances)
		}
		c := goccm.New(threads)
		var working []string
		for i := 0; i < len(instances); i++ {
			time.Sleep(time.Duration(cfg.Offset) * time.Millisecond)
			c.Wait()
			go func(i int) {
				err := instances[i].CheckToken()
				if err != 200 {
					color.Red("[%v] Токен недействителен %v", time.Now().Format("15:04:05"), instances[i].Token)
				} else {
					color.Green("[%v] Токен действителен %v", time.Now().Format("15:04:05"), instances[i].Token)
					working = append(working, instances[i].Token)
				}
				c.Done()
			}(i)
		}
		c.WaitAllDone()
		t := utilities.TruncateLines("tokens.txt", working)
		if t != nil {
			color.Red("[%v] Ошибка при усечении файла tokens.txt: %v", time.Now().Format("15:04:05"), t)
			ExitSafely()
			return
		}

		color.Green("[%v]Завершил все потоки", time.Now().Format("15:04:05"))

	case 8:
		// Leavs tokens from a server
		color.Cyan("Выйти из сервера")
		cfg, instances, err := getEverything()
		if err != nil {
			color.Red("Ошибка при получении необходимых данных %v", err)
			ExitSafely()

		}
		color.White("Введите количество потоков (0 для неограниченного количества): ")
		var threads int
		fmt.Scanln(&threads)
		if threads > len(instances) {
			threads = len(instances)
		}
		if threads == 0 {
			threads = len(instances)
		}
		color.White("Введите задержку между виходами: ")
		var delay int
		fmt.Scanln(&delay)
		color.White("Введите ID сервера: ")
		var serverid string
		fmt.Scanln(&serverid)
		c := goccm.New(threads)
		for i := 0; i < len(instances); i++ {
			time.Sleep(time.Duration(cfg.Offset) * time.Millisecond)
			c.Wait()
			go func(i int) {
				p := instances[i].Leave(serverid)
				if p == 0 {
					color.Red("[%v] Ошибка при выходе", time.Now().Format("15:04:05"))
				}
				if p == 200 || p == 204 {
					color.Green("[%v] Покинул сервер", time.Now().Format("15:04:05"))
				} else {
					color.Red("[%v] Ошибка при выходе", time.Now().Format("15:04:05"))
				}
				time.Sleep(time.Duration(delay) * time.Second)
				c.Done()
			}(i)
		}
		c.WaitAllDone()
		color.Green("[%v] Завершил все потоки", time.Now().Format("15:04:05"))
	case 9:

		color.Blue("Токен Onliner")
		_, instances, err := getEverything()
		if err != nil {
			color.Red("Ошибка при получении необходимых данных %v", err)
			ExitSafely()
		}
		var wg sync.WaitGroup
		wg.Add(len(instances))
		for i := 0; i < len(instances); i++ {
			go func(i int) {
				err := instances[i].StartWS()
				if err != nil {
					color.Red("[%v] Ошибка при открытии вебсокета: %v", time.Now().Format("15:04:05"), err)
				} else {
					color.Green("[%v] Открыт вебсокет %v", time.Now().Format("15:04:05"), instances[i].Token)
				}
				wg.Done()
			}(i)
		}
		wg.Wait()
		color.Green("[%v] Все Токены онлайн. Нажмите ENTER для отключения и продолжения программы", time.Now().Format("15:04:05"))
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		wg.Add(len(instances))
		for i := 0; i < len(instances); i++ {
			go func(i int) {
				instances[i].Ws.Close()
				wg.Done()
			}(i)
		}
		wg.Wait()
		color.Green("[%v] Все Токены offline", time.Now().Format("15:04:05"))

	case 10:
		color.Blue("Парсинг меню")
		cfg, _, err := getEverything()
		if err != nil {
			color.Red("Ошибка при получении необходимых данных %v", err)
		}
		color.White("1) Онлайн Парсинг\n2) Парсинг из Реакций\n3) Offline Парсинг")
		var options int
		fmt.Scanln(&options)
		if options == 1 {
			var token string
			color.White("Введите токен: ")
			fmt.Scanln(&token)
			var serverid string
			color.White("Введите  ID сервера: ")
			fmt.Scanln(&serverid)
			var channelid string
			color.White("Введите ID канала: ")
			fmt.Scanln(&channelid)

			Is := utilities.Instance{Token: token}
			t := 0
			for {
				if t >= 5 {
					color.Red("[%v] Не удалось подключиться к вебсокету после повторной попытки.", time.Now().Format("15:04:05"))
					break
				}
				err := Is.StartWS()
				if err != nil {
					color.Red("[%v] Ошибка при открытии вебсокета: %v", time.Now().Format("15:04:05"), err)
				} else {
					break
				}
				t++
			}

			color.Green("[%v] Открыт вебсокет %v", time.Now().Format("15:04:05"), Is.Token)

			i := 0
			for {
				err := utilities.Scrape(Is.Ws, serverid, channelid, i)
				if err != nil {
					color.Red("[%v] Ошибка при парсинге: %v", time.Now().Format("15:04:05"), err)
				}
				color.Green("[%v] Токен %v количество получено: %v", time.Now().Format("15:04:05"), Is.Token, len(Is.Ws.Members))
				if Is.Ws.Complete {
					break
				}
				i++
				time.Sleep(time.Duration(cfg.SleepSc) * time.Millisecond)
			}
			if Is.Ws != nil {
				Is.Ws.Close()
			}
			color.Green("[%v] Парсинг завершен. Получено %v пользователей", time.Now().Format("15:04:05"), len(Is.Ws.Members))
			clean := utilities.RemoveDuplicateStr(Is.Ws.Members)
			color.Green("[%v] Удалены дубликаты. Получено %v пользователей", time.Now().Format("15:04:05"), len(clean))
			color.Green("[%v] Записать в файл memberids.txt? (y/n)", time.Now().Format("15:04:05"))

			var write string
			fmt.Scanln(&write)
			if write == "y" {
				for k := 0; k < len(clean); k++ {
					err := utilities.WriteLines("memberids.txt", clean[k])
					if err != nil {
						color.Red("[%v] Ошибка при записи в файл memberids.txt: %v", time.Now().Format("15:04:05"), err)
					}
				}
				color.Green("[%v] Записано в файл memberids.txt", time.Now().Format("15:04:05"))
				err := WriteFile("парсинг/"+serverid+".txt", clean)
				if err != nil {
					color.Red("[%v] Ошибка при записи в файл: %v", time.Now().Format("15:04:05"), err)
				}
			}

		}
		if options == 2 {
			var token string
			color.White("Введите токен: ")
			fmt.Scanln(&token)
			var messageid string
			color.White("Введите ID сообщения: ")
			fmt.Scanln(&messageid)
			var channelid string
			color.White("Введите ID канала: ")
			fmt.Scanln(&channelid)
			color.White("1) Получить Emoji из сообщения\n2) Ввести Emoji вручную")
			var option int
			var send string
			fmt.Scanln(&option)
			var emoji string
			if option == 2 {
				color.White("Введите emoji [Для нативных эмодзи Discord просто скопируйте и вставьте emoji в формате unicode. Для пользовательских/нитро-эмодзи введите Имя:EmojiID именно в таком формате]: ")
				fmt.Scanln(&emoji)
				send = emoji
			} else {
				msg, err := utilities.GetRxn(channelid, messageid, token)
				if err != nil {
					fmt.Println(err)
				}
				color.White("Выберите Emoji")
				for i := 0; i < len(msg.Reactions); i++ {
					color.White("%v) %v %v", i, msg.Reactions[i].Emojis.Name, msg.Reactions[i].Count)
				}
				var index int
				fmt.Scanln(&index)
				if msg.Reactions[index].Emojis.ID == "" {
					send = msg.Reactions[index].Emojis.Name

				} else if msg.Reactions[index].Emojis.ID != "" {
					send = msg.Reactions[index].Emojis.Name + ":" + msg.Reactions[index].Emojis.ID
				}
			}

			var allUIDS []string
			var m string
			for {
				if len(allUIDS) == 0 {
					m = ""
				} else {
					m = allUIDS[len(allUIDS)-1]
				}
				rxn, err := utilities.GetReactions(channelid, messageid, token, send, m)
				if err != nil {
					fmt.Println(err)
					continue
				}
				if len(rxn) == 0 {
					break
				}
				fmt.Println(rxn)
				allUIDS = append(allUIDS, rxn...)

			}
			color.Green("[%v] Парсинг завершен. Полученно %v линии - Удаление дубликатов", time.Now().Format("15:04:05"), len(allUIDS))
			clean := utilities.RemoveDuplicateStr(allUIDS)
			color.Green("[%v] Записать в файл memberids.txt? (y/n)", time.Now().Format("15:04:05"))
			var write string
			fmt.Scanln(&write)
			if write == "y" {
				for k := 0; k < len(clean); k++ {
					err := utilities.WriteLines("memberids.txt", clean[k])
					if err != nil {
						color.Red("[%v] Ошибка при записи в файл memberids.txt: %v", time.Now().Format("15:04:05"), err)
					}
				}
				color.Green("[%v] Записано в файл memberids.txt", time.Now().Format("15:04:05"))
				err := WriteFile("парсинг/"+messageid+".txt", allUIDS)
				if err != nil {
					color.Red("[%v] Ошибка при записи в файл: %v", time.Now().Format("15:04:05"), err)
				}
			}
			fmt.Println("Готово")
		}
		if options == 3 {
			// Query Brute. This is a test function. Try using the compressed stream to appear legit.
			// Make a list of possible characters - Space can only come once, double spaces are counted as single ones and Name can't start from space. Queries are NOT case-sensitive.
			// Start from a character, check the returns. If it's less than 100, that query is complete and no need to go further down the rabbit hole.
			// If it's more than 100 or 100 and the last name starts from the query, pick the letter after our query and go down the rabbit hole.
			// Wait 0.5s (Or better, needs testing) Between scrapes and systematically connect and disconnect from websocket to avoid rate limiting.
			// Global var where members get appended (even repeats, will be cleared later) list of queries completed, list of queries left to complete and last query the instance searched to be in struct
			// Scan line for user input to stop at any point and proceed with the memberids scraped at hand.
			// Multiple instance support. Division of queries and hence completes in lesser time.
			// Might not need to worry about spaces at all as @ uses no spaces.
			// Starting Websocket(s) Appending to a slice. 1 for now, add more later.
			color.Cyan("Offline Парсинг")
			color.White("Эта функция намеренно замедляется с большими задержками. Пожалуйста, используйте несколько токенов и убедитесь, что они находятся на сервере перед началом работы, чтобы завершить ее быстро.")
			cfg, instances, err := getEverything()
			if err != nil {
				color.Red("[%v] Ошибка при получении конфигурации: %v", time.Now().Format("15:04:05"), err)
				ExitSafely()
			}
			var scraped []string
			// Input the number of tokens to be used
			color.Green("[%v] Сколько Токенов вы хотите использовать? У вас есть %v ", time.Now().Format("15:04:05"), len(instances))
			var numTokens int
			quit := make(chan bool)
			var allQueries []string
			fmt.Scanln(&numTokens)

			chars := " !\"#$%&'()*+,-./0123456789:;<=>?@[]^_`abcdefghijklmnopqrstuvwxyz{|}~"
			queriesLeft := make(chan string)
			var queriesCompleted []string

			for i := 0; i < len(chars); i++ {
				go func(i int) {
					queriesLeft <- string(chars[i])
				}(i)
			}

			if numTokens > len(instances) {
				color.Red("[%v] У вас только %v Токенов в файле tokens.txt Использую максимальное  количество токенов", time.Now().Format("15:04:05"), len(instances))
			} else if numTokens <= 0 {
				color.Red("[%v] Вы должны использовать как минимум 1 Токен", time.Now().Format("15:04:05"))
				ExitSafely()
			} else if numTokens <= len(instances) {
				color.Green("[%v] У вас есть %v токенов в файле tokens.txt Использование %v токенов", time.Now().Format("15:04:05"), len(instances), numTokens)
				instances = instances[:numTokens]
			} else {
				color.Red("[%v] Неверный ввод", time.Now().Format("15:04:05"))
			}

			color.Green("[%v] Введите  ID сервера", time.Now().Format("15:04:05"))
			var serverid string
			fmt.Scanln(&serverid)
			color.Green("[%v] Нажмите ENTER для запуска и остановки парсинга", time.Now().Format("15:04:05"))
			bufio.NewReader(os.Stdin).ReadBytes('\n')
			// Starting the instances as GOroutines
			for i := 0; i < len(instances); i++ {
				go func(i int) {
					instances[i].ScrapeCount = 0
					for {

						// Start websocket, reconnect if disconnected.
						if instances[i].ScrapeCount%5 == 0 || instances[i].LastCount%100 == 0 {
							if instances[i].Ws != nil {
								instances[i].Ws.Close()
							}
							time.Sleep(2 * time.Second)
							err := instances[i].StartWS()
							if err != nil {
								fmt.Println(err)
								continue
							}
							time.Sleep(2 * time.Second)

						}
						instances[i].ScrapeCount++

						// Get a query from the channel / Await for close response
						select {
						case <-quit:
							return
						default:
							query := <-queriesLeft
							allQueries = append(allQueries, query)
							if instances[i].Ws == nil {
								continue
							}
							if instances[i].Ws.Conn == nil {
								continue
							}
							err := utilities.ScrapeOffline(instances[i].Ws, serverid, query)
							if err != nil {
								color.Red("[%v] %v Ошибка при парсинге: %v", time.Now().Format("15:04:05"), instances[i].Token, err)
								go func() {
									queriesLeft <- query
								}()
								continue
							}

							memInfo := <-instances[i].Ws.OfflineScrape
							queriesCompleted = append(queriesCompleted, query)
							var MemberInfo utilities.Event
							err = json.Unmarshal(memInfo, &MemberInfo)
							if err != nil {
								color.Red("[%v] Error while unmarshalling: %v", time.Now().Format("15:04:05"), err)
								queriesLeft <- query
								continue
							}

							if len(MemberInfo.Data.Members) == 0 {
								instances[i].LastCount = -1
								continue
							}
							instances[i].LastCount = len(MemberInfo.Data.Members)
							for _, member := range MemberInfo.Data.Members {
								// Avoiding Duplicates
								if !utilities.Contains(scraped, member.User.ID) {
									scraped = append(scraped, member.User.ID)
								}
							}
							color.Green("[%v] Токен %v Запрос %v Получено %v [+%v]", time.Now().Format("15:04:05"), instances[i].Token, query, len(scraped), len(MemberInfo.Data.Members))

							for i := 0; i < len(MemberInfo.Data.Members); i++ {
								id := MemberInfo.Data.Members[i].User.ID
								err := utilities.WriteLines("memberids.txt", id)
								if err != nil {
									color.Red("[%v] Ошибка при записи в файл: %v", time.Now().Format("15:04:05"), err)
									continue
								}
							}
							if len(MemberInfo.Data.Members) < 100 {
								time.Sleep(time.Duration(cfg.SleepSc) * time.Millisecond)
								continue
							}
							lastName := MemberInfo.Data.Members[len(MemberInfo.Data.Members)-1].User.Username

							nextQueries := findNextQueries(query, lastName, queriesCompleted, chars)
							for i := 0; i < len(nextQueries); i++ {
								go func(i int) {
									queriesLeft <- nextQueries[i]
								}(i)
							}

						}

					}
				}(i)
			}

			bufio.NewReader(os.Stdin).ReadBytes('\n')
			color.Green("[%v] Остановка всех потоков", time.Now().Format("15:04:05"))
			for i := 0; i < len(instances); i++ {
				go func() {
					quit <- true
				}()
			}

			color.Green("[%v] Парсинг завершен. %v участников получено.", time.Now().Format("15:04:05"), len(scraped))
			color.Green("Хотите ли вы снова записать файл? (y/n) [Это удалит ранее существовавшие идентификаторы из файла memberids.txt]")
			var choice string
			fmt.Scanln(&choice)
			if choice == "y" || choice == "Y" {
				clean := utilities.RemoveDuplicateStr(scraped)
				err := utilities.TruncateLines("memberids.txt", clean)
				if err != nil {
					color.Red("[%v] Error while truncating file: %v", time.Now().Format("15:04:05"), err)
				}
				err = WriteFile("парсинг/"+serverid, clean)
				if err != nil {
					color.Red("[%v] Ошибка при записи в файл: %v", time.Now().Format("15:04:05"), err)
				}
			}

		}
	case 11:
		color.Blue("Смена имени")
		_, instances, err := getEverything()
		if err != nil {
			color.Red("[%v] Ошибка при получении необходимых данных: %v", time.Now().Format("15:04:05"), err)
		}
		for i := 0; i < len(instances); i++ {
			if !reg.MatchString(instances[i].Token) {
				color.Red("[%v] Программа смены имени требует токены в формате email:pass:token, возможно, токены неправильно отформатированы", time.Now().Format("15:04:05"))
				continue
			}
			fullz := instances[i].Token
			instances[i].Token = strings.Split(fullz, ":")[2]
			instances[i].Password = strings.Split(fullz, ":")[1]
		}
		color.Red("Примечание: Фотографии профиля изменяются случайным образом из файла.")
		users, err := utilities.ReadLines("names.txt")
		if err != nil {
			color.Red("[%v] При чтение файла ошибка names.txt: %v", time.Now().Format("15:04:05"), err)
			ExitSafely()
		}
		color.Green("[%v] Введите количество потоков: ", time.Now().Format("15:04:05"))

		var threads int
		fmt.Scanln(&threads)
		if threads > len(instances) {
			threads = len(instances)
		}

		c := goccm.New(threads)
		for i := 0; i < len(instances); i++ {
			c.Wait()
			go func(i int) {
				r, err := instances[i].NameChanger(users[rand.Intn(len(users))])
				if err != nil {
					color.Red("[%v] %v При изменении ИМЯ ошибка: %v", time.Now().Format("15:04:05"), instances[i].Token, err)
					return
				}
				body, err := utilities.ReadBody(r)
				if err != nil {
					fmt.Println(err)
				}
				if r.StatusCode == 200 || r.StatusCode == 204 {
					color.Green("[%v] %v ИМЯ успешно изменено", time.Now().Format("15:04:05"), instances[i].Token)
				} else {
					color.Red("[%v] %v При изменении ИМЯ ошибка: %v %v", time.Now().Format("15:04:05"), instances[i].Token, r.Status, string(body))
				}
				c.Done()
			}(i)
		}
		c.WaitAllDone()
		color.Green("[%v] Все сделано", time.Now().Format("15:04:05"))

	case 12:
		color.Blue("аватарка профиля изменен ")
		_, instances, err := getEverything()
		if err != nil {
			color.Red("[%v] Ошибка при получении необходимых данных: %v", time.Now().Format("15:04:05"), err)
		}
		color.Red("ПРИМЕЧАНИЕ: Поддерживаются только PNG и JPEG/JPG. Картинки профиля меняются случайным образом из папки. Используйте формат PNG для получения более быстрых результатов.")
		color.White("Загрузка аватаров...")
		ex, err := os.Executable()
		if err != nil {
			color.Red("Не смог найти Exe")
			ExitSafely()
		}
		ex = filepath.ToSlash(ex)
		path := path.Join(path.Dir(ex) + "/input/pfps")

		images, err := utilities.GetFiles(path)
		if err != nil {
			color.Red("Не удалось найти изображения в папке PFPs")
			ExitSafely()
		}
		color.Green("%v найденные файлы", len(images))
		var avatars []string

		for i := 0; i < len(images); i++ {
			av, err := utilities.EncodeImg(images[i])
			if err != nil {
				color.Red("Не удалось закодировать изображение")
				continue
			}
			avatars = append(avatars, av)
		}
		color.Green("%v  аватары загруженные", len(avatars))
		color.Green("[%v] Введите количество потоков: ", time.Now().Format("15:04:05"))
		var threads int
		fmt.Scanln(&threads)
		if threads > len(instances) {
			threads = len(instances)
		}

		c := goccm.New(threads)
		for i := 0; i < len(instances); i++ {
			c.Wait()

			go func(i int) {
				r, err := instances[i].AvatarChanger(avatars[rand.Intn(len(avatars))])
				if err != nil {
					color.Red("[%v] %v Ошибка при смене аватара: %v", time.Now().Format("15:04:05"), instances[i].Token, err)
				} else {
					if r.StatusCode == 204 || r.StatusCode == 200 {
						color.Green("[%v] %v Аватар успешно изменен", time.Now().Format("15:04:05"), instances[i].Token)
					} else {
						color.Red("[%v] %v Ошибка при смене аватара: %v", time.Now().Format("15:04:05"), instances[i].Token, r.StatusCode)
					}
				}

				c.Done()
			}(i)
		}
		c.WaitAllDone()
		color.Green("[%v] Все сделано", time.Now().Format("15:04:05"))
	case 13:
		color.White("Проверьте, находятся ли ваши токены на сервере")
		_, instances, err := getEverything()
		if err != nil {
			color.Red("[%v] Ошибка при получении необходимых данных: %v", time.Now().Format("15:04:05"), err)
			ExitSafely()
		}
		var serverid string
		var inServer []string
		color.Green("[%v] Введите ID сервера: ", time.Now().Format("15:04:05"))
		fmt.Scanln(&serverid)
		var wg sync.WaitGroup
		wg.Add(len(instances))
		for i := 0; i < len(instances); i++ {
			go func(i int) {
				defer wg.Done()
				r, err := instances[i].ServerCheck(serverid)
				if err != nil {
					color.Red("[%v] %v Ошибка при проверке сервера: %v", time.Now().Format("15:04:05"), instances[i].Token, err)
				} else {
					if r == 200 || r == 204 {
						color.Green("[%v] %v находится на сервере %v ", time.Now().Format("15:04:05"), instances[i].Token, serverid)
						inServer = append(inServer, instances[i].Token)
					} else if r == 429 {
						color.Green("[%v] %v ограничена по лимитах", time.Now().Format("15:04:05"), instances[i].Token)
					} else if r == 400 {
						color.Red("[%v] Плохой запрос - Неверный ID сервера", time.Now().Format("15:04:05"))
					} else {
						color.Red("[%v] %v не находится на сервере [%v] [%v]", time.Now().Format("15:04:05"), instances[i].Token, serverid, r)
					}
				}
			}(i)
		}
		wg.Wait()
		color.Green("[%v] Все готово. Вы хотите сохранить только токены на сервере в файл tokens.txt? (y/n)", time.Now().Format("15:04:05"))
		var save string
		fmt.Scanln(&save)
		if save == "y" || save == "Y" {
			err := utilities.TruncateLines("tokens.txt", inServer)
			if err != nil {
				color.Red("[%v]Ошибка при сохранении Токенов: %v", time.Now().Format("15:04:05"), err)
			} else {
				color.Green("[%v] Токены сохраняются в файле tokens.txt", time.Now().Format("15:04:05"))
			}
		}

	
	case 14:
		color.Blue("Made  by @M_rd2 for Private.")
	case 15:
		// Exit without error
		os.Exit(0)

	}
	time.Sleep(1 * time.Second)
	Options()

}

type jsonResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func getEverything() (utilities.Config, []utilities.Instance, error) {
	var cfg utilities.Config
	var instances []utilities.Instance
	var err error
	var tokens []string
	var proxies []string
	var proxy string

	// Load config
	cfg, err = utilities.GetConfig()
	if err != nil {
		return cfg, instances, err
	}
	if cfg.Proxy != "" && os.Getenv("HTTPS_PROXY") == "" {
		os.Setenv("HTTPS_PROXY", "http://"+cfg.Proxy)
	}
	if cfg.CaptchaAPI == "" {
		color.Red("[!] You're not using a Captcha API, some functionality like invite joining might be unavailable")
	}
	if !utilities.Contains(CaptchaServices, cfg.CaptchaAPI) {
		color.Red("[!] Captcha API %v is not supported. Please use one of the following: %v", cfg.CaptchaAPI, CaptchaServices)
		cfg.CaptchaAPI = ""
	}

	// Load instances
	tokens, err = utilities.ReadLines("tokens.txt")
	if err != nil {
		return cfg, instances, err
	}
	if len(tokens) == 0 {
		return cfg, instances, fmt.Errorf("no tokens found in tokens.txt")
	}
	if cfg.ProxyFromFile {
		proxies, err = utilities.ReadLines("proxies.txt")
		if err != nil {
			return cfg, instances, err
		}
		if len(proxies) == 0 {
			return cfg, instances, fmt.Errorf("no proxies found in proxies.txt")
		}
	}
	for i := 0; i < len(tokens); i++ {
		if cfg.ProxyFromFile {
			proxy = proxies[rand.Intn(len(proxies))]
		} else {
			proxy = ""
		}
		client, err := initClient(proxy, cfg)
		if err != nil {
			return cfg, instances, fmt.Errorf("couldn't initialize client: %v", err)
		}
		// proxy is put in struct only to be used by gateway. If proxy for gateway is disabled, it will be empty
		if !cfg.GatewayProxy {
			proxy = ""
		}
		instances = append(instances, utilities.Instance{Client: client, Token: tokens[i], Proxy: proxy, Config: cfg})
	}
	if len(instances) == 0 {
		color.Red("[!] You may be using 0 tokens")
	}
	var empty utilities.Config
	if cfg == empty {
		color.Red("[!] You may be using a malformed config.json")
	}
	return cfg, instances, nil

}

func setMessages(instances []utilities.Instance, messages []utilities.Message) error {
	var err error
	if len(messages) == 0 {
		messages, err = utilities.GetMessage()
		if err != nil {
			return err
		}
		if len(messages) == 0 {
			return fmt.Errorf("no messages found in messages.txt")
		}
		for i := 0; i < len(instances); i++ {
			instances[i].Messages = messages
		}
	} else {
		for i := 0; i < len(instances); i++ {
			instances[i].Messages = messages
		}
	}

	return nil
}

// Append items from slice to file
func Append(filename string, items []string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, item := range items {
		if _, err = file.WriteString(item + "\n"); err != nil {
			return err
		}
	}

	return nil
}

// Truncate items from slice to file
func Truncate(filename string, items []string) error {
	file, err := os.OpenFile(filename, os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, item := range items {
		if _, err = file.WriteString(item + "\n"); err != nil {
			return err
		}
	}

	return nil
}

// Write line to file
func WriteLine(filename string, line string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = file.WriteString(line + "\n"); err != nil {
		return err
	}

	return nil
}

// Create a New file and add items from a slice or append to it if it already exists
func WriteFile(filename string, items []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, item := range items {
		if _, err = file.WriteString(item + "\n"); err != nil {
			return err
		}
	}

	return nil
}

func initClient(proxy string, cfg utilities.Config) (*http.Client, error) {
	// If proxy is empty, return a default client (if proxy from file is false)
	if proxy == "" {
		return http.DefaultClient, nil
	}
	if !strings.Contains(proxy, "http://") {
		proxy = "http://" + proxy
	}
	// Error while converting proxy string to url.url would result in default client being returned
	proxyURL, err := url.Parse(proxy)
	if err != nil {
		return http.DefaultClient, err
	}
	// Creating a client and modifying the transport.

	Client := &http.Client{
		Timeout: time.Second * time.Duration(cfg.Timeout),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion:         tls.VersionTLS12,
				CipherSuites:       []uint16{0x1301, 0x1303, 0x1302, 0xc02b, 0xc02f, 0xcca9, 0xcca8, 0xc02c, 0xc030, 0xc00a, 0xc009, 0xc013, 0xc014, 0x009c, 0x009d, 0x002f, 0x0035},
				InsecureSkipVerify: true,
				CurvePreferences:   []tls.CurveID{tls.CurveID(0x001d), tls.CurveID(0x0017), tls.CurveID(0x0018), tls.CurveID(0x0019), tls.CurveID(0x0100), tls.CurveID(0x0101)},
			},
			DisableKeepAlives: cfg.DisableKL,
			ForceAttemptHTTP2: true,
			Proxy:             http.ProxyURL(proxyURL),
		},
	}
	return Client, nil

}

func ExitSafely() {
	color.Red("\nPress ENTER to EXIT")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	os.Exit(0)
}

const logo = "\r\n\r\n\u2588\u2588\u2588\u2588\u2588\u2588\u2557 \u2588\u2588\u2588\u2557   \u2588\u2588\u2588\u2557\u2588\u2588\u2588\u2588\u2588\u2588\u2557  \u2588\u2588\u2588\u2588\u2588\u2588\u2557  \u2588\u2588\u2588\u2588\u2588\u2588\u2557 \r\n\u2588\u2588\u2554\u2550\u2550\u2588\u2588\u2557\u2588\u2588\u2588\u2588\u2557 \u2588\u2588\u2588\u2588\u2551\u2588\u2588\u2554\u2550\u2550\u2588\u2588\u2557\u2588\u2588\u2554\u2550\u2550\u2550\u2550\u255D \u2588\u2588\u2554\u2550\u2550\u2550\u2588\u2588\u2557\r\n\u2588\u2588\u2551  \u2588\u2588\u2551\u2588\u2588\u2554\u2588\u2588\u2588\u2588\u2554\u2588\u2588\u2551\u2588\u2588\u2551  \u2588\u2588\u2551\u2588\u2588\u2551  \u2588\u2588\u2588\u2557\u2588\u2588\u2551   \u2588\u2588\u2551\r\n\u2588\u2588\u2551  \u2588\u2588\u2551\u2588\u2588\u2551\u255A\u2588\u2588\u2554\u255D\u2588\u2588\u2551\u2588\u2588\u2551  \u2588\u2588\u2551\u2588\u2588\u2551   \u2588\u2588\u2551\u2588\u2588\u2551   \u2588\u2588\u2551\r\n\u2588\u2588\u2588\u2588\u2588\u2588\u2554\u255D\u2588\u2588\u2551 \u255A\u2550\u255D \u2588\u2588\u2551\u2588\u2588\u2588\u2588\u2588\u2588\u2554\u255D\u255A\u2588\u2588\u2588\u2588\u2588\u2588\u2554\u255D\u255A\u2588\u2588\u2588\u2588\u2588\u2588\u2554\u255D\r\n\u255A\u2550\u2550\u2550\u2550\u2550\u255D \u255A\u2550\u255D     \u255A\u2550\u255D\u255A\u2550\u2550\u2550\u2550\u2550\u255D  \u255A\u2550\u2550\u2550\u2550\u2550\u255D  \u255A\u2550\u2550\u2550\u2550\u2550\u255D \r\nDISCORD MASS DM GO V1.7.9\n"

func findNextQueries(query string, lastName string, completedQueries []string, chars string) []string {
	if query == "" {
		color.Red("[%v] Query is empty", time.Now().Format("15:04:05"))
		return nil
	}
	lastName = strings.ToLower(lastName)
	indexQuery := strings.Index(lastName, query)
	if indexQuery == -1 {
		return nil
	}
	wantedCharIndex := indexQuery + len(query)
	if wantedCharIndex >= len(lastName) {

		return nil
	}
	wantedChar := lastName[wantedCharIndex]
	queryIndexDone := strings.Index(chars, string(wantedChar))
	if queryIndexDone == -1 {

		return nil
	}

	var nextQueries []string
	for j := queryIndexDone; j < len(chars); j++ {
		newQuery := query + string(chars[j])
		if !utilities.Contains(completedQueries, newQuery) && !strings.Contains(newQuery, "  ") && string(newQuery[0]) != "" {
			nextQueries = append(nextQueries, newQuery)
		}
	}
	return nextQueries
}

var CaptchaServices []string
