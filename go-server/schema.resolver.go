// schema.resolvers.go

func (r *mutationResolver) CreateMovie(ctx context.Context, input model.NewMovie) (*model.Movie, error) {
	movie := model.Movie{
	Title: input.Title,
	URL: input.URL,
	}
	
	_, err := r.DB.Model(&movie).Insert()
	if err != nil {
	return nil, fmt.Errorf("error inserting new movie: %v", err)
	}
	
	return &movie, nil
	}