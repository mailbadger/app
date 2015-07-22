<?php

namespace newsletters\Entities;

use Illuminate\Database\Eloquent\Model;
use Prettus\Repository\Contracts\Transformable;
use Prettus\Repository\Traits\TransformableTrait;

class Lists extends Model implements Transformable
{
    use TransformableTrait;

    protected $table = 'lists';

    protected $fillable = [
        'name',
        'total_subscribers',
    ];

    public function subscribers()
    {
        return $this->belongsToMany('newsletters\Entities\Subscriber', 'subscribers_lists', 'list_id',
            'subscriber_id')->withTimestamps();
    }
}
