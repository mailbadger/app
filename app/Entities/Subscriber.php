<?php

namespace newsletters\Entities;

use Illuminate\Database\Eloquent\Model;
use Prettus\Repository\Contracts\Transformable;
use Prettus\Repository\Traits\TransformableTrait;

class Subscriber extends Model implements Transformable
{
    use TransformableTrait;

    protected $table = 'subscribers';

    protected $fillable = [
        'email',
        'name',
    ];

    public function lists()
    {
        return $this->belongsToMany('newsletters\Entities\Lists', 'subscribers_lists', 'subscriber_id',
            'list_id')->withTimestamps();
    }

    public function fields()
    {
        return $this->belongsToMany('newsletters\Entities\Field', 'subscribers_fields', 'subscriber_id',
            'field_id')
            ->withPivot('value')
            ->withTimestamps();
    }
}
