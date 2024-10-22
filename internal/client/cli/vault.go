package cli

import (
	"context"
	"errors"
	"fmt"

	"github.com/pinbrain/gophkeeper/internal/model"
	"github.com/spf13/cobra"
)

func (c *CLI) GetDataCmd(ctx context.Context) *cobra.Command {
	var id string
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Загрузка данных из хранилища",
		Long:  "Загрузить данные из хранилища по id",
		RunE: func(_ *cobra.Command, _ []string) error {
			dataType, res, err := c.service.GetData(ctx, id)
			if err != nil {
				return err
			}
			switch dataType {
			case model.Password:
				item, ok := res.(*model.PasswordItem)
				if !ok {
					return errors.New("некорректные данные о пароле")
				}
				fmt.Printf(
					"Ресурс: %s; Логин: %s; Пароль: %s\nКомментарий: %s\n",
					item.Meta.Resource,
					item.Meta.Login,
					item.Data,
					item.Meta.Comment,
				)
				return nil

			case model.BankCard:
				item, ok := res.(*model.BankCardItem)
				if !ok {
					return errors.New("некорректные данные о банковской карте")
				}
				fmt.Printf(
					"Банк: %s\nДержатель: %s; Номер: %s; Действует до: %d-%d; csv: %s\nКомментарий: %s\n",
					item.Meta.Bank,
					item.Data.Holder,
					item.Data.Number,
					item.Data.ValidMonth,
					item.Data.ValidYear,
					item.Data.CSV,
					item.Meta.Comment,
				)
				return nil

			case model.Text:
				item, ok := res.(*model.TextItem)
				if !ok {
					return errors.New("некорректные данные о тексте")
				}
				fmt.Printf(
					"Название: %s\nТекст: %s\nКомментарий: %s\n",
					item.Meta.Name, item.Data, item.Meta.Comment,
				)
				return nil

			case model.File:
				item, ok := res.(*model.FileItem)
				if !ok {
					return errors.New("некорректные данные о файле")
				}
				fmt.Printf(
					"Файл '%s' успешно загружен.\n",
					item.Meta.Name,
				)
				return nil
			}
			return fmt.Errorf("неизвестный тип данных: %s", dataType)
		},
	}
	cmd.Flags().StringVar(&id, "id", "", "id данных для загрузки")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func (c *CLI) GetAllByTypeCmd(ctx context.Context) *cobra.Command {
	var dataType string
	cmd := &cobra.Command{
		Use:   "getall",
		Short: "Получить список данных",
		Long:  "Получить список хранящихся данных определенного типа",
		RunE: func(_ *cobra.Command, _ []string) error {
			res, err := c.service.GetAllByType(ctx, model.DataType(dataType))
			if err != nil {
				return err
			}
			if len(res) == 0 {
				fmt.Println("Данных нет")
				return nil
			}
			switch dataType {
			case string(model.Password):
				for _, item := range res {
					meta, ok := item.Meta.(*model.PasswordMeta)
					if !ok {
						return errors.New("не удалось получить мета данные о пароле")
					}
					fmt.Printf(
						"id: %s; Ресурс: %s; Логин: %s; Комментарий: %s\n",
						item.ID, meta.Resource, meta.Login, meta.Comment,
					)
				}

			case string(model.Text):
				for _, item := range res {
					meta, ok := item.Meta.(*model.TextMeta)
					if !ok {
						return errors.New("не удалось получить мета данные о тексте")
					}
					fmt.Printf(
						"id: %s; Имя: %s; Комментарий: %s\n",
						item.ID, meta.Name, meta.Comment,
					)
				}

			case string(model.BankCard):
				for _, item := range res {
					meta, ok := item.Meta.(*model.BankCardMeta)
					if !ok {
						return errors.New("не удалось получить мета данные о банковской карте")
					}
					fmt.Printf(
						"id: %s; Банк: %s; Комментарий: %s\n",
						item.ID, meta.Bank, meta.Comment,
					)
				}

			case string(model.File):
				for _, item := range res {
					meta, ok := item.Meta.(*model.FileMeta)
					if !ok {
						return errors.New("не удалось получить мета данные о файле")
					}
					fmt.Printf(
						"id: %s; Имя: %s; Расширение: %s; Комментарий: %s\n",
						item.ID, meta.Name, meta.Extension, meta.Comment,
					)
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&dataType, "type", "t", "", "тип данных для вывода списка")
	_ = cmd.MarkFlagRequired("type")
	return cmd
}

func (c *CLI) DeleteDataCmd(ctx context.Context) *cobra.Command {
	var id string
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Удалить",
		Long:  "Удалить данные по id",
		RunE: func(_ *cobra.Command, _ []string) error {
			err := c.service.DeleteData(ctx, id)
			if err != nil {
				return err
			}
			fmt.Println("Данные успешно удалены")
			return nil
		},
	}
	cmd.Flags().StringVar(&id, "id", "", "id данных для удаления")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func (c *CLI) AddDataCmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Добавить",
		Long:  "Загрузить данные в хранилище",
	}

	cmd.AddCommand(
		c.AddPasswordCmd(ctx),
		c.AddTextCmd(ctx),
		c.AddBankCardCmd(ctx),
		c.AddFileCmd(ctx),
	)

	return cmd
}

func (c *CLI) AddPasswordCmd(ctx context.Context) *cobra.Command {
	var password, login, resource, comment string
	cmd := &cobra.Command{
		Use:   "password",
		Short: "Добавить пароль",
		Long:  "Добавить в хранилище новый пароль",
		RunE: func(_ *cobra.Command, _ []string) error {
			err := c.service.AddPassword(ctx, password, model.PasswordMeta{
				Resource: resource,
				Login:    login,
				Comment:  comment,
			})
			if err != nil {
				return err
			}
			fmt.Println("Пароль успешно добавлен в хранилище")
			return nil
		},
	}
	cmd.Flags().StringVarP(&password, "password", "p", "", "пароль для хранения")
	_ = cmd.MarkFlagRequired("password")
	cmd.Flags().StringVarP(&login, "login", "l", "", "логин от ресурса")
	_ = cmd.MarkFlagRequired("login")
	cmd.Flags().StringVarP(&resource, "resource", "r", "", "название ресурса")
	_ = cmd.MarkFlagRequired("resource")
	cmd.Flags().StringVarP(&comment, "comment", "c", "", "комментарий")
	return cmd
}

func (c *CLI) AddTextCmd(ctx context.Context) *cobra.Command {
	var data, name, comment string
	cmd := &cobra.Command{
		Use:   "text",
		Short: "Добавить текст",
		Long:  "Добавить в хранилище текстовую информацию",
		RunE: func(_ *cobra.Command, _ []string) error {
			err := c.service.AddText(ctx, data, model.TextMeta{
				Name:    name,
				Comment: comment,
			})
			if err != nil {
				return err
			}
			fmt.Println("Текст успешно добавлен в хранилище")
			return nil
		},
	}
	cmd.Flags().StringVarP(&data, "text", "t", "", "текстовые данные")
	_ = cmd.MarkFlagRequired("text")
	cmd.Flags().StringVarP(&name, "name", "n", "", "название текста")
	_ = cmd.MarkFlagRequired("name")
	cmd.Flags().StringVarP(&comment, "comment", "c", "", "комментарий")
	return cmd
}

func (c *CLI) AddBankCardCmd(ctx context.Context) *cobra.Command {
	cardData := model.BankCardData{}
	var bankName, comment string
	cmd := &cobra.Command{
		Use:   "bcard",
		Short: "Добавить карту",
		Long:  "Добавить в хранилище данные банковской карты",
		RunE: func(_ *cobra.Command, _ []string) error {
			err := c.service.AddBankCard(ctx, cardData, model.BankCardMeta{
				Bank:    bankName,
				Comment: comment,
			})
			if err != nil {
				return err
			}
			fmt.Println("Банковская карта успешно добавлена в хранилище")
			return nil
		},
	}
	cmd.Flags().StringVarP(&cardData.Holder, "owner", "o", "", "держатель")
	_ = cmd.MarkFlagRequired("owner")
	cmd.Flags().StringVarP(&cardData.Number, "number", "n", "", "номер карты")
	_ = cmd.MarkFlagRequired("number")
	cmd.Flags().StringVarP(&cardData.CSV, "csv", "s", "", "csv код")
	_ = cmd.MarkFlagRequired("csv")
	cmd.Flags().IntVarP(&cardData.ValidMonth, "month", "m", 0, "месяц действия до")
	_ = cmd.MarkFlagRequired("month")
	cmd.Flags().IntVarP(&cardData.ValidYear, "year", "y", 0, "год действия до")
	_ = cmd.MarkFlagRequired("year")
	cmd.Flags().StringVarP(&bankName, "bank", "b", "", "банк")
	_ = cmd.MarkFlagRequired("bank")
	cmd.Flags().StringVarP(&comment, "comment", "c", "", "комментарий")
	return cmd
}

func (c *CLI) AddFileCmd(ctx context.Context) *cobra.Command {
	var file, comment string
	cmd := &cobra.Command{
		Use:   "file",
		Short: "Добавить файл",
		Long:  "Добавить в хранилище файл",
		RunE: func(_ *cobra.Command, _ []string) error {
			err := c.service.AddFile(ctx, file, comment)
			if err != nil {
				return err
			}
			fmt.Println("Файл успешно добавлен в хранилище")
			return nil
		},
	}
	cmd.Flags().StringVarP(&file, "path", "p", "", "файл")
	_ = cmd.MarkFlagRequired("path")
	cmd.Flags().StringVarP(&comment, "comment", "c", "", "комментарий")
	return cmd
}
